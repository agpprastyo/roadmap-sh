// src/pages/admin/Admin.jsx

import {useEffect, useState} from 'react';
import {useAuth} from '../../contexts/AuthContext';
import {useNavigate} from 'react-router-dom';
import {format} from 'date-fns';
import ArticleDialog from "../../components/ArticleDialog.jsx"; // Recommended for date formatting

const Admin = () => {
    const {isLoggedIn, logout} = useAuth(); // Get logged in status and logout function
    const navigate = useNavigate();
    const [articles, setArticles] = useState([]);
    const [search, setSearch] = useState('');
    const [page, setPage] = useState(1);
    const [pageSize] = useState(10); // Static page size
    const [totalArticles, setTotalArticles] = useState(0); // Total articles for pagination

    // const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedArticle, setSelectedArticle] = useState(null);
    const [dialogMode, setDialogMode] = useState('create');

    // Method to handle create
    const handleCreateArticle = () => {
        setDialogMode('create');
        setSelectedArticle(null);
        setIsDialogOpen(true);
    };

    // Method to handle edit
    const handleEditArticle = (articleId) => {
        const article = articles.find(a => a.id === articleId);
        setSelectedArticle(article);
        setDialogMode('edit');
        setIsDialogOpen(true);
    };

    // Method to submit created or edited article
    const handleArticleSubmit = async (articleId, articleData) => {
        try {
            const url = articleId
                ? `http://localhost:4444/api/v1/edit/${articleId}`
                : 'http://localhost:4444/api/v1/create';

            const method = articleId ? 'PATCH' : 'POST';

            const response = await fetch(url, {
                method,
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(articleData)
            });

            if (response.ok) {
                // Close dialog and refresh articles
                setIsDialogOpen(false);
                await fetchArticles();

                // Optional: Show success message
                alert(articleId
                    ? 'Article updated successfully!'
                    : 'Article created successfully!'
                );
            } else {
                // Handle error
                const errorData = await response.json();
                console.error('Error submitting article:', errorData);
                alert(`Failed to ${articleId ? 'update' : 'create'} article`);
            }
        } catch (error) {
            console.error('Error submitting article:', error);
            alert('An error occurred while submitting the article');
        }
    };

    const handleRestoreArticle = async (articleId) => {
        try {
            const response = await fetch(`http://localhost:4444/api/v1/restore/${articleId}`, {
                method: 'POST',
                credentials: 'include'
            });

            if (response.ok) {
                // Refresh articles after successful restoration
                await fetchArticles();
                alert('Article restored successfully');
            } else {
                const errorData = await response.json();
                console.error('Error restoring article:', errorData);
                alert('Failed to restore article');
            }
        } catch (error) {
            console.error('Error restoring article:', error);
            alert('An error occurred while restoring the article');
        }
    };


    // New method to handle article creation
    // const createArticle = async (articleData) => {
    //     try {
    //         const response = await fetch('http://localhost:4444/api/v1/create', {
    //             method: 'POST',
    //             credentials: 'include',
    //             headers: {
    //                 'Content-Type': 'application/json',
    //             },
    //             body: JSON.stringify(articleData)
    //         });
    //
    //         if (response.ok) {
    //             // Close dialog and refresh articles
    //             setIsCreateDialogOpen(false);
    //             fetchArticles();
    //
    //             // Optional: Show success message
    //             alert('Article created successfully!');
    //         } else {
    //             // Handle error
    //             const errorData = await response.json();
    //             console.error('Error creating article:', errorData);
    //             alert('Failed to create article');
    //         }
    //     } catch (error) {
    //         console.error('Error creating article:', error);
    //         alert('An error occurred while creating the article');
    //     }
    // };


    // Fetch articles when login status, page or search term changes
    useEffect(() => {
        if (isLoggedIn) {
            fetchArticles().then(r => r);
        }
    }, [isLoggedIn, page, search]);

    // Fetch articles from the API
    const fetchArticles = async () => {
        try {
            const response = await fetch(`http://localhost:4444/api/v1/admin?page=${page}&search=${search}&page_size=${pageSize}&sort_by=title-asc`, {
                method: 'GET',
                credentials: 'include', // Ensure cookies are sent with the request
            });

            if (response.ok) {
                const data = await response.json();
                setArticles(data.articles || []);
                setTotalArticles(data.total || 0); // Set total number of articles for pagination
            } else {
                console.error('Error fetching articles:', response.statusText);
            }
        } catch (error) {
            console.error('Error fetching articles:', error);
        }
    };

    // Toggle publish status of an article
    const togglePublished = async (articleId, currentStatus) => {
        // Ensure currentStatus is a boolean
        const newStatus = !currentStatus; // Toggle the current status
        try {
            const response = await fetch(`http://localhost:4444/api/v1/edit/${articleId}`, {
                method: 'PATCH',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                    'Cookie': document.cookie // Assuming session cookies are already set
                },
                body: JSON.stringify({published: newStatus}) // Send the new status
            });

            if (response.ok) {
                // Re-fetch articles to get the updated list
                await fetchArticles();
            } else {
                console.error('Error updating article:', response.statusText);
            }
        } catch (error) {
            console.error('Error updating article:', error);
        }
    };

    // Handle searching articles
    const handleSearch = (e) => {
        e.preventDefault();
        setPage(1); // Reset to page 1 on new search
        fetchArticles();
    };

    // Handle logout
    const handleLogout = () => {
        logout(); // Call logout from context
        navigate('/login'); // Redirect to login page
    };

    // Handle page change
    const handlePageChange = (newPage) => {
        setPage(newPage);
    };

    const handleDeleteArticle = async (articleId) => {
        // Confirm deletion
        const confirmDelete = window.confirm('Are you sure you want to delete this article?');

        if (confirmDelete) {
            try {
                const response = await fetch(`http://localhost:4444/api/v1/delete/${articleId}`, {
                    method: 'DELETE',
                    credentials: 'include'
                });

                if (response.ok) {
                    // Refresh articles after successful deletion
                    fetchArticles();
                    alert('Article deleted successfully');
                } else {
                    const errorData = await response.json();
                    console.error('Error deleting article:', errorData);
                    alert('Failed to delete article');
                }
            } catch (error) {
                console.error('Error deleting article:', error);
                alert('An error occurred while deleting the article');
            }
        }
    };


    // Calculate total pages
    const totalPages = Math.ceil(totalArticles / pageSize);

    return (
        <div className="max-w-6xl mx-auto">
            <header className="bg-white shadow-md">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center py-6">
                        <div className="flex items-center">
                            <h1 className="text-3xl font-extrabold text-gray-900">
                                Admin Dashboard
                            </h1>
                        </div>
                        <div className="flex items-center space-x-4">
                            <button
                                onClick={handleCreateArticle}
                                className="bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600 transition-colors flex items-center space-x-2"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20"
                                     fill="currentColor">
                                    <path fillRule="evenodd"
                                          d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z"
                                          clipRule="evenodd"/>
                                </svg>
                                <span>Create Article</span>
                            </button>
                            <button
                                onClick={handleLogout}
                                className="bg-red-500 text-white px-4 py-2 rounded-md hover:bg-red-600 transition-colors flex items-center space-x-2"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20"
                                     fill="currentColor">
                                    <path fillRule="evenodd"
                                          d="M3 3a1 1 0 00-1 1v12a1 1 0 102 0V4a1 1 0 00-1-1zm10.293 1.293a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 01-1.414-1.414L14.586 11H7a1 1 0 110-2h7.586l-1.293-1.293a1 1 0 010-1.414z"
                                          clipRule="evenodd"/>
                                </svg>
                                <span>Logout</span>
                            </button>
                        </div>
                    </div>
                </div>
            </header>

            <main className="mx-auto sm:px-6 lg:px-1 py-8">
                <div className="bg-white shadow-md rounded-lg p-6">
                    <div className="flex items-center justify-center my-6">
                        <form onSubmit={handleSearch} className="flex items-center space-x-2 w-full max-w-md">
                            <div className="relative flex-grow">
                                <input
                                    type="text"
                                    placeholder="Search articles..."
                                    value={search}
                                    onChange={(e) => setSearch(e.target.value)}
                                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    className="h-5 w-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                                    viewBox="0 0 20 20"
                                    fill="currentColor"
                                >
                                    <path fillRule="evenodd"
                                          d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                                          clipRule="evenodd"/>
                                </svg>
                            </div>
                            <button
                                type="submit"
                                className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 transition-colors"
                            >
                                Search
                            </button>
                        </form>
                    </div>
                </div>
            </main>


            <table className="min-w-full bg-white border border-gray-300 mb-4 shadow-sm rounded-lg overflow-hidden">
                <thead className="bg-gray-100 border-b">
                <tr>
                    <th className="border px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">
                        #
                    </th>
                    <th className="border px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Title
                    </th>
                    <th className="border px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Created At
                    </th>
                    <th className="border px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Deleted At
                    </th>
                    <th className="border px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Published
                    </th>
                    <th className="border px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                    </th>
                </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                {articles.length === 0 ? (
                    <tr>
                        <td
                            className="px-4 py-4 text-center text-gray-500"
                            colSpan="6"
                        >
                            No articles found
                        </td>
                    </tr>
                ) : (
                    articles.map((article, index) => (
                        <tr
                            key={article.id}
                            className="hover:bg-gray-50 transition-colors"
                        >
                            <td className="border px-4 py-2 text-gray-600">
                                {index + 1}
                            </td>
                            <td className="border px-4 py-2 font-medium text-gray-900">
                                {article.title}
                            </td>
                            <td className="border px-4 py-2 text-gray-600">
                                {article.created_at
                                    ? format(new Date(article.created_at), 'PPpp')
                                    : 'N/A'}
                            </td>
                            <td className="border px-4 py-2 text-gray-600">
                                {article.delete_at  // Note: delete_at instead of deleted_at
                                    ? format(new Date(article.delete_at), 'PPpp')
                                    : 'N/A'}
                            </td>
                            <td className=" px-4 py-2 flex items-center justify-center">
                                <div
                                    className="relative inline-block"
                                    title={article.delete_at ? "Cannot modify deleted article" : ""}
                                >
                                    <button
                                        onClick={() => togglePublished(article.id, article.published)}
                                        className={`flex items-center space-x-2 ${
                                            article.delete_at && article.delete_at !== '0001-01-01T00:00:00Z'
                                                ? 'opacity-50 cursor-not-allowed'
                                                : ''
                                        }`}
                                        aria-pressed={article.published}
                                        aria-label={`${article.published ? 'Unpublish' : 'Publish'} article`}
                                        disabled={article.delete_at && article.delete_at !== '0001-01-01T00:00:00Z'}
                                    >
                                        <div
                                            className={`w-10 h-5 rounded-full relative transition-colors duration-200 ${
                                                article.published
                                                    ? 'bg-green-500'
                                                    : 'bg-red-500'
                                            } ${
                                                article.delete_at && article.delete_at !== '0001-01-01T00:00:00Z'
                                                    ? 'opacity-50'
                                                    : ''
                                            }`}
                                        >
                                            <div
                                                className={`absolute left-0 top-0 w-5 h-5 bg-white rounded-full shadow-md transform transition-transform duration-200 ${
                                                    article.published
                                                        ? 'translate-x-full'
                                                        : 'translate-x-0'
                                                }`}
                                            />
                                        </div>
                                        <span className={`text-xs font-medium ${
                                            article.published
                                                ? 'text-green-700'
                                                : 'text-red-700'
                                        } ${
                                            article.delete_at && article.delete_at !== '0001-01-01T00:00:00Z'
                                                ? 'opacity-50'
                                                : ''

                                        }`}>
                {article.published ? 'Published' : 'Unpublished'}
            </span>
                                    </button>
                                </div>
                            </td>

                            <td className="border px-4 py-2">
                                <div className="flex space-x-2">
                                    {/* If article is deleted (has delete_at), show restore button */}
                                    {article.delete_at && article.delete_at !== '0001-01-01T00:00:00Z' ? (
                                        <button
                                            onClick={() => handleRestoreArticle(article.id)}
                                            className="text-green-500 hover:text-green-700 flex items-center"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1"
                                                 viewBox="0 0 20 20" fill="currentColor">
                                                <path fillRule="evenodd"
                                                      d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z"
                                                      clipRule="evenodd"/>
                                            </svg>
                                            Restore
                                        </button>
                                    ) : (
                                        <>
                                            <button
                                                onClick={() => handleEditArticle(article.id)}
                                                className="text-blue-500 hover:text-blue-700 flex items-center"
                                            >
                                                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1"
                                                     viewBox="0 0 20 20" fill="currentColor">
                                                    <path
                                                        d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"/>
                                                </svg>
                                                Edit
                                            </button>
                                            <button
                                                onClick={() => handleDeleteArticle(article.id)}
                                                className="text-red-500 hover:text-red-700 flex items-center"
                                            >
                                                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1"
                                                     viewBox="0 0 20 20" fill="currentColor">
                                                    <path fillRule="evenodd"
                                                          d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                                                          clipRule="evenodd"/>
                                                </svg>
                                                Delete
                                            </button>
                                        </>
                                    )}
                                </div>
                            </td>

                        </tr>
                    ))
                )}
                </tbody>
            </table>


            {/* Pagination controls */}
            <div className="flex justify-start items-center gap-6">
                <button
                    onClick={() => handlePageChange(page > 1 ? page - 1 : 1)}
                    className={`py-1 px-4  hover:bg-blue-300 rounded ${page === 1 && 'opacity-50 cursor-not-allowed'}`}
                    disabled={page === 1}>
                    Previous
                </button>
                <span>Page {page} of {totalPages}</span>
                <button
                    onClick={() => handlePageChange(page < totalPages ? page + 1 : totalPages)}
                    className={`py-1 px-4 border bg-blue-600 text-white hover:bg-blue-300 rounded ${page === totalPages && 'opacity-50 cursor-not-allowed'}`}
                    disabled={page === totalPages}>
                    Next
                </button>
            </div>

            <ArticleDialog
                isOpen={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSubmit={handleArticleSubmit}
                initialData={selectedArticle}
                dialogMode={dialogMode}
            />
        </div>
    );
};

export default Admin;
