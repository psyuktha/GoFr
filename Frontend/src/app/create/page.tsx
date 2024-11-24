"use client"

import React, { useState } from 'react';
import { useSearchParams } from 'next/navigation';

const CreatePage = () => {
  const searchParams = useSearchParams();
  const platform = searchParams.get('platform');
  const [prompt, setPrompt] = useState('');
  const [response, setResponse] = useState(null);
  const [isConfirmed, setIsConfirmed] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const res = await fetch('https://api.sample.com/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ prompt }),
    });
    const data = await res.json();
    setResponse(data);
  };

  const handleFinalSubmit = async () => {
    const res = await fetch('https://api.sample.com/final-submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ prompt }),
    });
    const data = await res.json();
    // Handle final response
  };

  return (
    <div className="flex flex-col md:flex-row items-center justify-center min-h-screen py-2">
      <div className="w-full md:w-3/4 flex flex-col items-center">
        <h1 className="text-4xl font-bold text-black mb-4">{platform ? platform.charAt(0).toUpperCase() + platform.slice(1) : ''}</h1>
        <form className="w-full max-w-sm" onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="prompt">
              Prompt:
            </label>
            <textarea
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline h-32"
              id="prompt"
              name="prompt"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
            />
          </div>
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
            type="submit"
          >
            Send
          </button>
        </form>
        {response && (
          <div className="mt-4">
            <h2 className="text-2xl font-bold mb-2">Response</h2>
            <p>{response.message}</p>
            <div className="mt-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="confirm">
                Is this correct?
              </label>
              <button
                className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline mr-2"
                onClick={() => setIsConfirmed(true)}
              >
                Yes
              </button>
              <button
                className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                onClick={() => setIsConfirmed(false)}
              >
                No
              </button>
            </div>
            {isConfirmed && (
              <button
                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline mt-4"
                onClick={handleFinalSubmit}
              >
                Submit Final
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default CreatePage;