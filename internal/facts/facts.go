package facts

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// randomFact returns a random cat fact
//
// Taken from:
// https://github.com/vadimdemedes/cat-facts/blob/49dfacbe897b369f5403565b4d17614e459c468c/cat-facts.json
func randomFact() string {
	facts := []string{
		"Although it is known to be the tailless cat, the Manx can be born with a stub or a short tail",
		"Most cat litters contain four to six kittens",
		"On average, cats spend 2/3 of every day sleeping",
		"Blue-eyed cats have a high tendency to be deaf, but not all cats with blue eyes are deaf",
		"Researchers are unsure exactly how a cat purrs",
		"A cat almost never meows at another cat, mostly just humans",
		"Mohammed loved cats and reportedly his favorite cat, Muezza, was a tabby",
		"In homes with more than one cat, it is best to have cats of the opposite sex. They tend to be better housemates.",
		"Cats have over 100 sounds in their vocal repertoire, while dogs only have 10",
		"Cats would rather starve themselves than eat something they don't like. This means they will refuse an unpalatable -- but nutritionally complete -- food for a prolonged period",
		"The smallest pedigreed cat is a Singapura, which can weigh just 4 lbs",
		"Cats have a strong aversion to anything citrus",
		"Talk about Facetime: Cats greet one another by rubbing their noses together",
		"Black cats aren't an omen of ill fortune in all cultures. In the UK and Australia, spotting a black cat is good luck",
		"Most cats will eat 7 to 20 small meals a day. This interesting fact is brought to you by Nature's Recipe®",
		"One of Muhammad's companions was nicknamed Abu Hurairah, or Father of the Kitten, because he loved cats",
		"Outdoor cats' lifespan averages at about 3 to 5 years; indoor cats have lives that last 16 years or more",
		"Cats use their whiskers to measure openings, indicate mood and general navigation",
		"A cat's field of vision does not cover the area right under its nose",
		"Cats hate the water because their fur does not insulate well when it's wet",
		"During the Middle Ages, cats were associated with witchcraft",
		"The largest cat breed by mean weight is the Savannah, at 10kg",
		"A group of cats is called a clowder",
		"According to the Guinness World Records, the largest domestic cat litter totaled at 19 kittens, four of them stillborn",
		"A fingerprint is to a human as a nose is to a cat",
		"Genetically, cats' brains are more similar to that of a human than a dog's brain",
		"Landing on all fours is something typical to cats thanks to the help of their eyes and special balance organs in their inner ear. These tools help them straighten themselves in the air and land upright on the ground.",
		"In 1888, more than 300,000 mummified cats were found an Egyptian cemetery",
		"Many Egyptians worshipped the goddess Bast, who had a woman's body and a cat's head",
		"A cat cannot climb head first down a tree because its claws are curved the wrong way",
		"Some cats have survived falls of over 20 meters",
		"Twenty-five percent of cat owners use a blow drier on their cats after bathing",
		"Unlike dogs, cats do not have a sweet tooth",
		"Caution during Christmas: poinsettias may be festive, but they’re poisonous to cats",
		"When a family cat died in ancient Egypt, family members would mourn by shaving off their eyebrows",
		"Cats came to the Americas from Europe as pest controllers in the 1750s",
		"A cat usually has about 12 whiskers on each side of its face",
		"A cat's heart beats nearly twice as fast as a human heart",
		"According to the Association for Pet Obesity Prevention (APOP), about 50 million of our cats are overweight",
		"Female cats tend to be right pawed, while male cats are more often left pawed",
		"Cats who eat too much tuna can become addicted, which can actually cause a Vitamin E deficiency",
		"When a cat chases its prey, it keeps its head level",
		"A female cat is also known to be called a queen or a molly",
		"In one litter of kittens, there could be multiple father cats",
		"Cats can pick up on your tone of voice, so sweet-talking to your cat has more of an impact than you think",
		"Eating grass rids a cats' system of any fur and helps with digestion",
		"Cats make about 100 different sounds",
		"Two members of the cat family are distinct from all others: the clouded leopard and the cheetah",
		"Teeth of cats are sharper when they're kittens. After six months, they lose their needle-sharp milk teeth",
		"It is important to include fat in your cat's diet because they're unable to make the nutrient in their bodies on their own",
		"Many cat owners think their cats can read their minds",
		"Most cats give birth to a litter of between one and nine kittens",
		"Unlike most other cats, the Turkish Van breed has a water-resistant coat and enjoys being in water",
		"If your cat's eyes are closed, it's not necessarily because it's tired. A sign of closed eyes means your cat is happy or pleased",
		"The Snow Leopard, a variety of the California Spangled Cat, always has blue eyes",
		"A cat's brain is biologically more similar to a human brain than it is to a dog's",
		"A cat's eyesight is both better and worse than humans",
		"Rather than nine months, cats' pregnancies last about nine weeks",
		"A cat's meow is usually not directed at another cat, but at a human. To communicate with other cats, they will usually hiss, purr and spit.",
		"In Japan, cats are thought to have the power to turn into super spirits when they die",
		"In North America, cats are a more popular pet than dogs. Nearly 73 million cats and 63 million dogs are kept as household pets",
		"Around the world, cats take a break to nap —a catnap— 425 million times a day",
		"The smallest wildcat today is the Black-footed cat",
		"The color of York Chocolates becomes richer with age. Kittens are born with a lighter coat than the adults",
		"Because of widespread cat smuggling in ancient Egypt, the exportation of cats was a crime punishable by death",
		"A Japanese cat figurine called Maneki-Neko is believed to bring good luck",
		"There are more than 500 million domestic cats in the world",
		"The earliest ancestor of the modern cat lived about 30 million years ago",
		"Despite appearing like a wild cat, the Ocicat does not have an ounce of wild blood",
		"In multi-pet households, cats are able to get along especially well with dogs if they're introduced when the cat is under 6 months old and the dog is under one year old",
		"Want to call a hairball by its scientific name? Next time, say the word bezoar",
		"A cat can travel at a top speed of approximately 31 mph (49 km) over a short distance",
		"Cats have the skillset that makes them able to learn how to use a toilet",
		"Maine Coons are the most massive breed of house cats. They can weigh up to around 24 pounds",
		"Cats CAN be lefties and righties, just like us. More than forty percent of them are, leaving some ambidextrous",
		"Cats' rough tongues enable them to clean themselves efficiently and to lick clean an animal bone",
		"Smuggling a cat out of ancient Egypt was punishable by death",
		"Each side of a cat's face has about 12 whiskers",
		"Some cats can survive falls from as high up as 65 feet or more",
		"Most cats don't have eyelashes",
		"It has been said that the Ukrainian Levkoy has the appearance of a dog, due to the angles of its face",
		"Cats have 32 muscles that control the outer ear",
		"As temperatures rise, so do the number of cats. Cats are known to breed in warm weather, which leads many animal advocates worried about the plight of cats under Global Warming.",
		"Cats spend nearly 1/3 of their waking hours cleaning themselves",
		"A cat can reach up to five times its own height per jump",
		"The world's most fertile cat, whose name was Dusty, gave birth to 420 kittens in her lifetime",
		"The cat who holds the record for the longest non-fatal fall is Andy",
		"The Maine Coon is appropriately the official State cat of its namesake state",
		"Bobtails are known to have notably short tails -- about half or a third the size of the average cat",
		"Cats are extremely sensitive to vibrations",
		"Most kittens are born with blue eyes, which then turn color with age",
		"Cats actually have dreams, just like us. They start dreaming when they reach a week old",
		"The richest cat is Blackie who was left £15 million by his owner, Ben Rea",
		"Cat's back claws aren't as sharp as the claws on their front paws",
		"Cats sleep 16 hours of any given day",
		"A third of cats' time spent awake is usually spent cleaning themselves",
		"A cat's hearing is better than a dog's",
		"Most cats had short hair until about 100 years ago, when it became fashionable to own cats and experiment with breeding",
		"A cat's heart beats almost double the rate of a human heart, from 110 to 140 beats per minute",
		"A cat can jump up to five times its own height in a single bound",
		"Call them wide-eyes: cats are the mammals with the largest eyes",
		"A Selkirk slowly loses its naturally-born curly coat, but it grows again when the cat is around 8 months",
		"The two outer layers of a cat's hair are called, respectively, the guard hair and the awn hair",
		"Foods that should not be given to cats include onions, garlic, green tomatoes, raw potatoes, chocolate, grapes, and raisins",
		"A cat's jaw can't move sideways, so a cat can't chew large chunks of food",
		"Webbed feet on a cat? The Peterbald's got 'em! They make it easy for the cat to get a good grip on things with skill",
		"The Egyptian Mau is probably the oldest breed of cat",
		"Elvis Presley’s Chinese name is Mao Wong, or Cat King",
		"Cats show affection and mark their territory by rubbing on people. Glands on their face, tail and paws release a scent to make its mark",
		"Cats are the most popular pet in North American Cats are North America's most popular pets",
		"The biggest wildcat today is the Siberian Tiger",
		"Cats are unable to detect sweetness in anything they taste",
		"Collectively, kittens yawn about 200 million time per hour",
		"Cats have about 20,155 hairs per square centimeter",
		"If you killed a cat in the ages of Pharaoh, you could've been put to death",
		"The first cat show was organized in 1871 in London",
		"A cat has 230 bones in its body",
		"Today, cats are living twice as long as they did just 50 years ago",
		"Cats have the cognitive ability to sense a human's feelings and overall mood",
		"Approximately 40,000 people are bitten by cats in the U.S.",
		"A group of kittens is called a kindle, and clowder is a term that refers to a group of adult cats",
		"Cats have 24 more bones than humans",
		"The technical term for a cat's hairball is a bezoar",
		"Every year, nearly four million cats are eaten in Asia",
		"Perhaps the oldest cat breed on record is the Egyptian Mau, which is also the Egyptian language's word for cat",
		"When a household cat died in ancient Egypt, its owners showed their grief by shaving their eyebrows",
		"Cats prefer their food at room temperature—not too hot, not too cold",
		"Ragdoll cats live up to their name: they will literally go limp, with relaxed muscles, when lifted by a human",
		"Grown cats have 30 teeth",
		"Ancient Egyptians first adored cats for their finesse in killing rodents—as far back as 4,000 years ago",
		"Sir Isaac Newton, among his many achievements, invented the cat flap door",
		"Cats have a 5 toes on their front paws and 4 on each back paw",
		"Approximately 24 cat skins can make a coat",
		"Sometimes called the Canadian Hairless, the Sphynx is the first cat breed that has lasted this long—the breed has been around since 1966",
		"A cat's back is extremely flexible because it has up to 53 loosely fitting vertebrae",
		"According to the International Species Information Service, there are only three Marbled Cats still in existence worldwide.  One lives in the United States.",
	}
	n := rand.Int() % len(facts)
	return facts[n]
}

type completionRequest struct {
	User      string `json:"user"`
	MaxTokens int    `json:"max_tokens"`
	Prompt    string `json:"prompt"`
}

// GenerateFact generates a random fact using go-gpt3
func GenerateFact(id uint) (string, bool) {
	// Cheating with this instead of loading from static config
	secretKey := os.Getenv("CF_OPENAI_SECRET_KEY")
	if secretKey == "" {
		log.Println("CF_OPENAI_SECRET_KEY not set")
		return randomFact(), false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// convert id to string
	user := strconv.Itoa(int(id))

	completionURL := "https://api.openai.com/v1/engines/text-davinci-002/completions"
	completionRequest := completionRequest{
		User:      user,
		MaxTokens: 300,
		Prompt:    "write a wholesome story about cats or kittens without saying once upon a time",
	}

	jsonBody, err := json.Marshal(completionRequest)
	if err != nil {
		log.Println("Error marshalling completion request:", err)
		return randomFact(), false
	}

	// Create http req with context
	req, err := http.NewRequestWithContext(ctx, "POST", completionURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Error creating completion request:", err)
		return randomFact(), false
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error completing request:", err)
		return randomFact(), false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading completion response:", err)
		return randomFact(), false
	}

	type completionChoices struct {
		Text string `json:"text"`
	}

	type completionResponse struct {
		Choices []completionChoices `json:"choices"`
	}

	var response completionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error unmarshalling completion response:", err)
		return randomFact(), false
	}

	if len(response.Choices) == 0 {
		log.Println("No completion choices found")
		return randomFact(), false
	}

	return strings.TrimSpace(response.Choices[0].Text), true
}
