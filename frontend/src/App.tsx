import { ChangeEvent, FormEvent, useState } from 'react'
import './App.css'
import { QueryClient, QueryClientProvider, useQuery } from 'react-query'

const queryClient = new QueryClient()

interface Shipping {
  pack_size: number;
  shipping_pack_quantity: number;
}

interface Order {
  id: number;
  number_of_items: number;
  created_at: string;
  shipping: Shipping[];
}


async function fetchOrders() {
  var res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/orders`, { method: "GET" });
  if (!res.ok) throw new Error(`Response status: ${res.status}`);
  const json = await res.json();
  return json.data;
}


function App() {
  const [numItems, setNumItems] = useState(0);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const { data: orders, status } = useQuery<Order[]>("orders", fetchOrders);

  function handleChange(e: ChangeEvent<HTMLInputElement>) {
    setNumItems(Number(e.target.value));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setIsLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/orders`, {
        method: "POST",
        body: JSON.stringify({ number_of_items: numItems })
      })
      if (!res.ok) throw new Error(`Response status: ${res.status}`);
      alert("order have been created");
    } catch (e) {
      alert("could not create order");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <QueryClientProvider client={queryClient}>
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
                disabled={isLoading}
                className="inline-block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm disabled:bg-indigo-500 hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Add Order
              </button>
            </form>
          </div>
        </div>
        <div className="mt-8 flow-root">
          <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
              {status == "success" ?
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
                        Packaging
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {orders ? orders.map((order: Order) => (
                      <tr key={order.id}>
                        <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">
                          Order #{order.id}
                        </td>
                        <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{order.number_of_items}</td>
                        <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{new Date(order.created_at).toLocaleString()}</td>
                        <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                          {order.shipping.length ? order.shipping.map((s, key) => {
                            return <p key={key}>{s.shipping_pack_quantity} X {s.pack_size}</p>
                          }) : "-"}
                        </td>
                      </tr>
                    )) : ""}
                  </tbody>
                </table> : "No orders available"}
            </div>
          </div>
        </div>
      </div>
    </QueryClientProvider>
  )
}

export default App
