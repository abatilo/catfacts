import { useCallback, useState, Fragment } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { ExclamationIcon, CheckIcon } from "@heroicons/react/outline";

export default function HomePage() {
  const [open, setOpen] = useState(false);
  const [registerSuccessful, setRegisterSuccessful] = useState(true);
  const [phoneNumber, setPhoneNumberInput] = useState("");

  const submit = useCallback(
    async (event) => {
      event.preventDefault();
      const resp = await fetch("/api/register", {
        method: "POST",
        body: JSON.stringify({ phoneNumber }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (resp.ok) {
        setRegisterSuccessful(true);
      } else {
        setRegisterSuccessful(false);
      }
      setOpen(true);
    },
    [phoneNumber]
  );

  return (
    <>
      <Transition.Root show={open} as={Fragment}>
        <Dialog
          as="div"
          static
          className="fixed inset-0 z-10 overflow-y-auto"
          open={open}
          onClose={() => setOpen(false)}
        >
          <div className="flex items-end justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0"
              enterTo="opacity-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100"
              leaveTo="opacity-0"
            >
              <Dialog.Overlay className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
            </Transition.Child>

            {/* This element is to trick the browser into centering the modal contents. */}
            <span
              className="hidden sm:inline-block sm:align-middle sm:h-screen"
              aria-hidden="true"
            >
              &#8203;
            </span>
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              enterTo="opacity-100 translate-y-0 sm:scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0 sm:scale-100"
              leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            >
              <div className="inline-block px-4 pt-5 pb-4 overflow-hidden text-left align-bottom bg-white rounded-lg shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-sm sm:w-full sm:p-6">
                <div>
                  {registerSuccessful ? (
                    <div className="flex items-center justify-center w-12 h-12 mx-auto bg-green-100 rounded-full">
                      <CheckIcon
                        className="w-6 h-6 text-green-600"
                        aria-hidden="true"
                      />
                    </div>
                  ) : (
                    <div className="flex items-center justify-center w-12 h-12 mx-auto bg-red-100 rounded-full">
                      <ExclamationIcon
                        className="w-6 h-6 text-red-600"
                        aria-hidden="true"
                      />
                    </div>
                  )}
                  <div className="mt-3 text-center sm:mt-5">
                    <Dialog.Title
                      as="h3"
                      className="text-lg font-medium text-gray-900 leading-6"
                    >
                      {registerSuccessful
                        ? "You've registered"
                        : "Failed to register"}
                    </Dialog.Title>
                    <div className="mt-2">
                      <p className="text-sm text-gray-500">
                        {registerSuccessful
                          ? "You should receive a confirmation text message. You're on your way to receiving CatFacts!"
                          : "Your number couldn't be validated and stored. Could you double check the phone number? We only support US based phone numbers at this time."}
                      </p>
                    </div>
                  </div>
                </div>
                <div className="mt-5 sm:mt-6">
                  <button
                    type="button"
                    className="inline-flex justify-center w-full px-4 py-2 text-base font-medium text-white bg-yellow-600 border border-transparent rounded-md shadow-sm hover:bg-yellow-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500 sm:text-sm"
                    onClick={() => setOpen(false)}
                  >
                    Okay
                  </button>
                </div>
              </div>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>

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
                    Would you like AI to teach you more about them?
                  </span>
                </span>
              </h1>
              <p className="mt-3 text-base text-gray-500 sm:mt-5 sm:text-xl lg:text-lg xl:text-xl">
                For better or for worse, you've stumbled upon my little Twilio
                testing site. If you input your phone number below, you will
                start receiving text messages once a day with facts about cats!

                These facts are artificially generated by an{" "}
                <a href="https://openai.com/api"
                  className="font-medium text-gray-900 underline">
                  OpenAI model
                </a>.
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
    </>
  );
}
