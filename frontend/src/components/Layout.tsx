import { Menu, X } from "lucide-react";
import { useState } from "react";
import { NavLink, Outlet } from "react-router-dom";
import Footer from "./Footer";

export default function Layout() {
    const [sidebarOpen, setSidebarOpen] = useState(false);

    const toggleSidebar = () => setSidebarOpen(!sidebarOpen);

    return (
        <div className="min-h-screen flex bg-gray-100 text-gray-100">
            {/* Sidebar */}
            <aside
                className={`bg-gray-900 text-white w-64 space-y-6 py-7 px-2 absolute inset-y-0 left-0 transform ${sidebarOpen ? "translate-x-0" : "-translate-x-full"
                    } md:relative md:translate-x-0 transition duration-200 ease-in-out z-50`}
            >
                <div className="flex gap-4">
                    <button onClick={toggleSidebar} className="text-gray-100 text-2xl md:hidden">
                        <X />
                    </button>
                    <h1 className="text-2xl font-bold text-center">Algo-Arena</h1>
                </div>
                <nav className="flex flex-col space-y-2 mt-6">
                    {/* <NavItem to="/dashboard" label="Dashboard" /> */}
                    <NavItem to="/problems" label="Problems" />
                    {/* <NavItem to="/leaderboard" label="Leaderboard" /> */}
                    {/* <NavItem to="/profile" label="Profile" /> */}
                </nav>
            </aside>

            {/* Content */}
            <div className="flex-1 flex flex-col">
                {/* Top Bar */}
                <header className="bg-gray-900 shadow-md flex items-center justify-Start gap-4 px-4 py-3 md:hidden">
                    <button onClick={toggleSidebar} className="text-gray-100 text-2xl">
                        <Menu />
                    </button>
                    <h2 className="text-xl font-semibold text-gray-100">Algo-Arena</h2>
                </header>

                {/* Main Content */}
                <main className="flex-1 p-4 md:p-6 bg-gray-900 overflow-auto">
                    <Outlet />
                </main>

                <Footer />
            </div>
        </div>
    );
}

function NavItem({ to, label }: { to: string; label: string }) {
    return (
        <NavLink
            to={to}
            className={({ isActive }) =>
                `block py-2 px-4 rounded transition ${isActive ? "bg-indigo-600 font-semibold" : "hover:bg-gray-700"
                }`
            }
        >
            {label}
        </NavLink>
    );
}
