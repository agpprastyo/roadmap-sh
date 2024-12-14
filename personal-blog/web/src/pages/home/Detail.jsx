import {useState, useEffect} from 'react';
import {useParams, useNavigate} from 'react-router-dom';
import axios from 'axios';
import DOMPurify from 'dompurify';

const Detail = () => {
    const {id} = useParams();
    const navigate = useNavigate();
    const [article, setArticle] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchArticleDetail = async () => {
            try {
                setIsLoading(true);
                const response = await axios.get(`http://localhost:4444/api/v1/article/${id}`);
                setArticle(response.data);
            } catch (err) {
                setError(err);
                console.error('Error fetching article details:', err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchArticleDetail();
    }, [id]);

    const formatDate = (dateString) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const handleGoBack = () => {
        navigate(-1);
    };

    // Render loading state
    if (isLoading) {
        return (
            <div className="container mx-auto px-4 py-8 text-center">
                <div className="text-xl">Loading article...</div>
            </div>
        );
    }

    // Render error state
    if (error) {
        return (
            <div className="container mx-auto px-4 py-8 text-center">
                <div className="text-xl text-red-500">
                    Error loading article: {error.message}
                </div>
                <button
                    onClick={handleGoBack}
                    className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
                >
                    Go Back
                </button>
            </div>
        );
    }

    // Render article details
    return (
        <div className="max-w-4xl mx-auto py-8 px-4">
            {/* Back button */}
            <div className="mb-6">
                <button
                    onClick={handleGoBack}
                    className="flex items-center text-blue-600 hover:text-blue-800 transition"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-5 w-5 mr-2"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                    >
                        <path
                            fillRule="evenodd"
                            d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z"
                            clipRule="evenodd"
                        />
                    </svg>
                    Back to Articles
                </button>
            </div>

            <article className="bg-white shadow-lg rounded-lg overflow-hidden">
                <header className="px-6 py-4 border-b">
                    <h1 className="text-4xl font-bold text-gray-800 mb-2">{article.title}</h1>
                    <div className="text-gray-600 text-sm">
                        <span>Published: {formatDate(article.created_at)}</span>
                        {article.updated_at && (
                            <span className="ml-4">
                                Updated: {formatDate(article.updated_at)}
                            </span>
                        )}
                    </div>
                </header>

                <div
                    className="prose max-w-none px-6 py-4"
                    dangerouslySetInnerHTML={{
                        __html: DOMPurify.sanitize(article.content)
                    }}
                />
            </article>
        </div>
    );
};

export default Detail;
