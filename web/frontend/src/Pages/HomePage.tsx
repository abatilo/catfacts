import { useCallback, useState } from "react";

export default function HomePage() {
  const [phoneNumber, setPhoneNumberInput] = useState("");

  const submit = useCallback(
    async (event) => {
      event.preventDefault();
      await fetch("/api/register", {
        method: "POST",
        body: JSON.stringify({ phoneNumber }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      // setOpen(true);
    },
    [phoneNumber]
  );

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="pb-32 bg-yellow-600">
        <header className="py-10">
          <div className="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
            <h1 className="text-3xl font-bold text-white">
              Subscribe to CatFacts
            </h1>
          </div>
        </header>
      </div>

      <main className="-mt-32">
        <div className="px-4 pb-12 mx-auto max-w-7xl sm:px-6 lg:px-8">
          <div className="px-5 py-6 bg-white rounded-lg shadow sm:px-6">
            <h1>
              <span className="block mt-1 text-4xl font-extrabold tracking-tight sm:text-5xl xl:text-6xl">
                <span className="block text-gray-700">Do you love cats?</span>
                <span className="block text-gray-700">
                  Would you like to know more about them?
                </span>
              </span>
            </h1>
            <p className="mt-3 text-base text-gray-500 sm:mt-5 sm:text-xl lg:text-lg xl:text-xl">
              For better or for worse, you've stumbled upon my little Twilio
              testing site. If you input your phone number below, you will start
              receiving text messages once a day with facts about cats!
            </p>
            <div className="mt-8 sm:max-w-lg sm:mx-auto sm:text-center lg:text-left lg:mx-0">
              <p className="text-base font-medium text-gray-900">
                Sign up to start getting your facts
              </p>
              <form onSubmit={submit} className="mt-3 sm:flex">
                <label htmlFor="phone" className="sr-only">
                  Phone number
                </label>
                <input
                  autoFocus
                  type="tel"
                  name="phone"
                  id="phone"
                  className="block w-full py-3 text-base placeholder-gray-500 border-gray-300 rounded-md shadow-sm focus:ring-yellow-500 focus:border-yellow-500 sm:flex-1"
                  placeholder="Enter your phone number"
                  onChange={(e) => setPhoneNumberInput(e.target.value)}
                />
                <button
                  type="submit"
                  className="w-full px-6 py-3 mt-3 text-base font-medium text-white bg-gray-800 border border-transparent rounded-md shadow-sm hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500 sm:mt-0 sm:ml-3 sm:flex-shrink-0 sm:inline-flex sm:items-center sm:w-auto"
                >
                  Subscribe me
                </button>
              </form>
            </div>
            <p className="mt-3 text-sm text-gray-500">
              If you would like to, check out{" "}
              <a
                href="https://www.github.com/abatilo/catfacts"
                className="font-medium text-gray-900 underline"
              >
                the code on my GitHub
              </a>
              .
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}
