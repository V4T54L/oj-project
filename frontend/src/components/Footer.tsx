import { Copyright } from "lucide-react"

const Footer = () => {
    return (
        <footer className="text-center py-6 bg-gray-900 border-t border-gray-700 text-sm flex justify-center items-center gap-2">
           Copyright <Copyright size={12} /> {new Date().getFullYear()} Algo-Arena. All rights reserved.
        </footer>
    )
}

export default Footer