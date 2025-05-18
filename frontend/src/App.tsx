import { BrowserRouter, Route, Routes } from "react-router-dom"
import LandingPage from "./pages/LandingPage"
import NotFoundPage from "./pages/NotFoundPage"
import Layout from "./components/Layout"
import ProblemListPage from "./pages/ProblemsList"
import ViewProblemPage from "./pages/ViewProblem"

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/" element={<Layout />} >
          <Route path="problems" element={<ProblemListPage />} />
          <Route path="problems/:id" element={<ViewProblemPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Route>
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App