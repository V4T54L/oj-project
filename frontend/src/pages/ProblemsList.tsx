import { useState, useMemo, useEffect } from "react";
import type { ProblemInfo } from "../types"
import { NavLink } from "react-router-dom";
import { getProblems } from "../api/problems";

const ITEMS_PER_PAGE = 10;

type SortOption =
    | "title-asc"
    | "title-desc"
    | "acceptance-asc"
    | "acceptance-desc"
    | "difficulty-asc"
    | "difficulty-desc";

export default function ProblemListPage() {
    const [currentPage, setCurrentPage] = useState(1);
    const [searchTerm, setSearchTerm] = useState("");
    const [sortOption, setSortOption] = useState<SortOption>("title-asc");
    const [problems, setProblems] = useState<ProblemInfo[]>([]);

    const filteredProblems = useMemo(() => {
        const filtered = problems.filter((problem) =>
            problem.title.toLowerCase().includes(searchTerm.toLowerCase())
        );

        switch (sortOption) {
            case "title-asc":
                filtered.sort((a, b) => a.title.localeCompare(b.title));
                break;
            case "title-desc":
                filtered.sort((a, b) => b.title.localeCompare(a.title));
                break;
            case "acceptance-asc":
                filtered.sort((a, b) => a.acceptance_rate - b.acceptance_rate);
                break;
            case "acceptance-desc":
                filtered.sort((a, b) => b.acceptance_rate - a.acceptance_rate);
                break;
            case "difficulty-asc":
                filtered.sort((a, b) => a.difficulty.id - b.difficulty.id);
                break;
            case "difficulty-desc":
                filtered.sort((a, b) => b.difficulty.id - a.difficulty.id);
                break;
        }

        return filtered;
    }, [searchTerm, sortOption]);

    const totalPages = Math.ceil(filteredProblems.length / ITEMS_PER_PAGE);
    const currentProblems = filteredProblems.slice(
        (currentPage - 1) * ITEMS_PER_PAGE,
        currentPage * ITEMS_PER_PAGE
    );

    const handlePageChange = (page: number) => {
        if (page >= 1 && page <= totalPages) setCurrentPage(page);
    };

    useEffect(() => {
        const fetchProblems = async () => {
            try {
                const fetchedProblems = await getProblems();
                setProblems(fetchedProblems)
            } catch (error: any) {
                console.log("Error fetching the problems:", error)
            }
        }

        fetchProblems();
    }, [])

    return (
        <div className="bg-gray-900 text-white min-h-screen px-4 py-8 md:px-16">
            <h1 className="text-4xl font-bold mb-8 text-center bg-gradient-to-r from-indigo-500 to-purple-600 bg-clip-text text-transparent">
                Problems
            </h1>

            {/* Controls */}
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-8">
                <input
                    type="text"
                    placeholder="Search by title..."
                    className="bg-gray-800 border border-gray-600 text-white rounded px-4 py-2 w-full md:max-w-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
                    value={searchTerm}
                    onChange={(e) => {
                        setSearchTerm(e.target.value);
                        setCurrentPage(1);
                    }}
                />

                <select
                    value={sortOption}
                    onChange={(e) => setSortOption(e.target.value as SortOption)}
                    className="bg-gray-800 border border-gray-600 text-white rounded px-4 py-2 w-full md:w-64 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                    <option value="title-asc">Title A-Z</option>
                    <option value="title-desc">Title Z-A</option>
                    <option value="acceptance-desc">Acceptance ↑</option>
                    <option value="acceptance-asc">Acceptance ↓</option>
                    <option value="difficulty-asc">Difficulty Easy → Hard</option>
                    <option value="difficulty-desc">Difficulty Hard → Easy</option>
                </select>
            </div>

            {/* Table */}
            <div className="overflow-x-auto rounded-lg shadow border border-gray-700">
                <table className="min-w-full table-auto bg-gray-800">
                    <thead className="bg-gray-700 text-left text-sm font-semibold text-white">
                        <tr>
                            <th className="px-4 py-3">Title</th>
                            <th className="px-4 py-3">Difficulty</th>
                            <th className="px-4 py-3">Tags</th>
                            <th className="px-4 py-3">Acceptance</th>
                            <th className="px-4 py-3">Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {currentProblems.map((problem) => (
                            <tr key={problem.id} className="border-b border-gray-700 hover:bg-gray-700">
                                <td className="px-4 py-3">
                                    <NavLink
                                        to={`/problems/${problem.id}`}
                                        className="text-indigo-400 hover:underline"
                                    >
                                        {problem.title}
                                    </NavLink>
                                </td>
                                <td className="px-4 py-3">{problem.difficulty.name}</td>
                                <td className="px-4 py-3">
                                    <div className="flex flex-wrap gap-1">
                                        {problem.tags.map((tag) => (
                                            <span
                                                key={tag.id}
                                                className="bg-purple-700 text-white text-xs px-2 py-0.5 rounded"
                                            >
                                                {tag.name}
                                            </span>
                                        ))}
                                    </div>
                                </td>
                                <td className="px-4 py-3">{problem.acceptance_rate.toFixed(2)}%</td>
                                <td className="px-4 py-3">
                                    {problem.is_solved ? (
                                        <span className="text-green-400 font-medium">Solved</span>
                                    ) : (
                                        <span className="text-red-400 font-medium">Unsolved</span>
                                    )}
                                </td>
                            </tr>
                        ))}
                        {currentProblems.length === 0 && (
                            <tr>
                                <td colSpan={5} className="text-center py-6 text-gray-400">
                                    No problems found.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>

            {/* Pagination */}
            <div className="flex justify-center items-center gap-4 mt-8">
                <button
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 1}
                    className="px-4 py-2 border border-gray-600 rounded text-white hover:bg-gray-700 disabled:opacity-50"
                >
                    Prev
                </button>
                <span>
                    Page <strong>{currentPage}</strong> of {totalPages}
                </span>
                <button
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={currentPage === totalPages}
                    className="px-4 py-2 border border-gray-600 rounded text-white hover:bg-gray-700 disabled:opacity-50"
                >
                    Next
                </button>
            </div>
        </div>
    );
}
