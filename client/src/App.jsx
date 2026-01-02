import { Outlet } from 'react-router-dom';
import Navbar from './components/Navbar';

export default function App() {
  return (
    <div className="min-h-screen bg-base-200 flex flex-col">
      <Navbar />

      <main className="flex-1 flex items-center justify-center">
        <Outlet />
      </main>
    </div>
  );
}
