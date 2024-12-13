import { ChangeEvent, FormEvent, useEffect, useState } from 'react'
import './App.css'

function App() {
  const [numItems, setNumItems] = useState(0);
  const [orders, setOrders] = useState([
    { id: 1, quantity: 1, created_at: '2024-04-09', updated_at: '2024-04-09' },
  ]);

  function handleChange(e: ChangeEvent<HTMLInputElement>) {
    setNumItems(Number(e.target.value));
  }

  function handleSubmit(e: FormEvent) {
    e.preventDefault();

  }

  async function fetchOrders() {
    const BASE_URL = "http://localhost:8000";
    try {
      var res = await fetch(`${BASE_URL}/shipping`, { method: "GET" })
      if (!res.ok) throw new Error(`Response status: ${res.status}`);
      const json = await res.json();
      console.log(json);
    } catch (e) {
      alert("error fetching orders");
    }
  }

  useEffect(() => { fetchOrders() }, [])
  return (
    <>
      <div className="px-4 sm:px-6 lg:px-8 py-4">
        <div className="sm:flex sm:items-center">
          <div className="sm:flex-auto">
            <h1 className="text-base font-semibold text-gray-900">Orders</h1>
            <p className="mt-2 text-sm text-gray-700">
              A list of all the orders you have created
            </p>
          </div>
          <div className="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
            <form className='space-x-3' onSubmit={handleSubmit}>
              <div className='inline-block'>
                <input
                  onChange={handleChange}
                  id="num_items"
                  name="num_items"
                  type="number"
                  placeholder="0"
                  aria-label="Number of Items"
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                />
              </div>
              <button
                type="submit"
                className="inline-block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Add Order
              </button>
            </form>
          </div>
        </div>
        <div className="mt-8 flow-root">
          <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
              <table className="min-w-full divide-y divide-gray-300">
                <thead>
                  <tr>
                    <th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-0">
                      Order
                    </th>
                    <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Quantity
                    </th>
                    <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Date Created
                    </th>
                    <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Date Updated
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {orders.map((order) => (
                    <tr key={order.id}>
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">
                        Order #{order.id}
                      </td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{order.quantity}</td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{order.created_at}</td>
                      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{order.updated_at}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>


    </>
  )
}

export default App
