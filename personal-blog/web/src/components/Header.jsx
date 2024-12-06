import { useState } from "react";
import { Link } from "react-router-dom";
import { RoutePaths } from "../general/RoutePaths.jsx";

const Header = ({ onSearch, onSort }) => {
    const [searchTerm, setSearchTerm] = useState('');

    const handleSearch = () => {
        onSearch(searchTerm);
    };

    const handleKeyDown = (e) => {
        // Trigger search when Enter key is pressed
        if (e.key === 'Enter') {
            handleSearch();
        }
    };

    return (
        <header className="bg-blue-500 text-white p-6">
            <div className="max-w-6xl mx-auto flex justify-between items-center">
                <Link to={RoutePaths.HOME}>
                    <h1 className="text-2xl font-bold
                        hover:text-blue-200
                        transition-all
                        duration-300
                        transform
                        hover:scale-105">
                        AGP Blog
                    </h1>
                </Link>
                <div className="flex items-center space-x-2">
                    <input
                        type="text"
                        placeholder="Search articles..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        onKeyDown={handleKeyDown} // Add this line
                        className="px-2 py-1 text-black"
                    />
                    <button
                        onClick={handleSearch}
                        className="bg-white text-blue-500 px-3 py-1 rounded"
                    >
                        Search
                    </button>
                    <select
                        onChange={(e) => onSort(e.target.value)}
                        className="px-2 py-1 text-black"
                    >
                        <option value="">Sort By</option>
                        <option value="title-asc">Title (A-Z)</option>
                        <option value="title-desc">Title (Z-A)</option>
                        <option value="created_at-desc">Newest</option>
                        <option value="created_at-asc">Oldest</option>
                    </select>
                </div>
            </div>
        </header>
    );
};

export default Header;
