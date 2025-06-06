import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom"
import AuthPage from "./pages/AuthPage"
import ProblemListPage from "./pages/ProblemListPage";
import ProblemDetailPage from "./pages/ProblemDetailPage";
import NotFoundPage from "./pages/NotFoundPage";
import ContestListPage from "./pages/ContestListPage";
import ContestDetailPage from "./pages/ContestDetailPage";
import Navbar from "./components/Navbar";
import { useEffect, useState } from "react";
import { getCurrentUser } from "./api/endpoints";
import type { User } from "./types";
import UserProfile from "./pages/UserProfile";

const App = () => {
  const [user, setUser] = useState<User>()

  const getCurrentUserProfile = async () => {
    try {
      const res = await getCurrentUser();
      setUser(res.data)
    } catch (err) {
      console.log("Error: ", err)
      // setServerError(err.response?.data?.message || 'Something went wrong.');
    }
  };

  useEffect(() => {
    getCurrentUserProfile();
  }, [])

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50 text-gray-900">
        <Navbar user={user} />
        <Routes>
          <Route path="/" element={<Navigate to="/problems" />} />
          <Route path="/auth" element={<AuthPage />} />
          <Route path="/problems" element={<ProblemListPage />} />
          <Route path="/problem/:slug" element={<ProblemDetailPage />} />
          <Route path="/contests" element={<ContestListPage />} />
          <Route path="/contest/:contestId" element={<ContestDetailPage />} />
          {/* <Route path="/add-problem" element={<AddProblemPage />} /> */}
          <Route path="/profile/:username" element={<UserProfile />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App