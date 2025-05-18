import { NavLink } from "react-router-dom";

export default function NotFoundPage() {
    return (
        <div className="bg-gray-900 text-white min-h-screen flex flex-col items-center justify-center text-center px-4">
            <h1 className="text-6xl font-bold mb-4">404</h1>
            <h2 className="text-2xl md:text-3xl font-semibold mb-4">Page Not Found</h2>
            <p className="text-gray-400 mb-8">
                The page you’re looking for doesn’t exist or has been moved.
            </p>
            <NavLink
                to="/"
                className="bg-white text-indigo-700 px-6 py-3 rounded-full font-semibold hover:bg-gray-200 transition"
            >
                Go to Homepage
            </NavLink>
        </div>
    );
}
