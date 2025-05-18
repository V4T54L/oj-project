import { type ReactNode } from "react";
import Footer from "../components/Footer";
import { NavLink } from "react-router-dom";

export default function LandingPage() {
    return (
        <div className="bg-gray-900 text-white min-h-screen flex flex-col">
            {/* Hero Section */}
            <header className="bg-gradient-to-br from-indigo-600 to-purple-700 p-8 md:p-16 text-center flex-1">
                <h1 className="text-4xl md:text-6xl font-bold mb-4">Welcome to Algo-Arena</h1>
                <p className="text-xl md:text-2xl text-gray-200 mb-6">
                    Sharpen your coding skills. Conquer every algorithm. Rule the verse.
                </p>
                <div className="flex justify-center gap-4">
                    <NavLink
                        to="/problems"
                        className="bg-white text-indigo-700 px-6 py-3 rounded-full font-semibold hover:bg-gray-200 transition"
                    >
                        Get Started
                    </NavLink>
                    <NavLink
                        to="/problems"
                        className="border border-white px-6 py-3 rounded-full hover:bg-white hover:text-purple-700 transition"
                    >
                        Explore Problems
                    </NavLink>
                </div>
            </header>

            {/* Features Section */}
            <section className="flex flex-col md:flex-row justify-around text-center py-12 px-4 md:px-16 bg-gray-800">
                <FeatureCard
                    title="Daily Challenges"
                    description="Solve new problems every day to stay sharp and consistent."
                    icon="🧠"
                />
                <FeatureCard
                    title="Leaderboard"
                    description="Compete with coders globally and climb the AlgoVerse rankings."
                    icon="🏆"
                />
                <FeatureCard
                    title="Detailed Solutions"
                    description="Get in-depth explanations and community discussions."
                    icon="📘"
                />
            </section>

            {/* Footer */}
            <Footer />
        </div>
    );
}

function FeatureCard({ title, description, icon }: { title: string, description: string, icon: ReactNode }) {
    return (
        <div className="bg-gray-700 p-6 rounded-lg max-w-sm mb-6 md:mb-0">
            <div className="text-4xl mb-4">{icon}</div>
            <h3 className="text-xl font-semibold mb-2">{title}</h3>
            <p className="text-gray-300">{description}</p>
        </div>
    );
}
