import {useState, useEffect, useCallback} from 'react';
import axios from 'axios';
import Header from "../../components/Header.jsx";
import ArticleCard from "../../components/ArticleCard.jsx";
import Pagination from "../../components/Pagination.jsx";

const Home = () => {
    const [articles, setArticles] = useState([]);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [searchTerm, setSearchTerm] = useState('');
    const [sortBy, setSortBy] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);


    const fetchArticles = useCallback(async () => {
        setIsLoading(true);
        setError(null);
        try {
            // Create dynamic params object based on sortBy
            const sortParams = {};
            switch (sortBy) {
                case 'title-asc':
                    sortParams.title_asc = true;
                    break;
                case 'title-desc':
                    sortParams.title_desc = true;
                    break;
                case 'created_at-desc':
                    sortParams.created_at_desc = true;
                    break;
                case 'created_at-asc':
                    sortParams.created_at_asc = true;
                    break;
                default:
                    break;
            }

            const response = await axios.get('http://localhost:4444/api/v1/articles', {
                params: {
                    page,
                    search: searchTerm,
                    page_size: 10,
                    ...sortParams
                }
            });

            // Comprehensive null and undefined checks
            const articlesData = response?.data?.articles ?? [];
            const total = response?.data?.total ?? 0;

            setArticles(Array.isArray(articlesData) ? articlesData : []);
            setTotalPages(Math.ceil(total / 10));
        } catch (error) {
            console.error('Error fetching articles:', error);
            setError(error);
            // Reset articles and total pages in case of error
            setArticles([]);
            setTotalPages(0);
        } finally {
            setIsLoading(false);
        }
    }, [page, searchTerm, sortBy]);

    useEffect(() => {
        fetchArticles();
    }, [fetchArticles]);

    return (
        <div>
            <Header
                onSearch={(term) => {
                    setSearchTerm(term);
                    setPage(1); // Reset to first page on new search
                }}
                onSort={setSortBy}
            />
            <div className="max-w-6xl mx-auto py-4">
                {isLoading && (
                    <div className="text-center text-gray-500">
                        Loading articles...
                    </div>
                )}

                {error && (
                    <div className="text-center text-red-500">
                        Error: {error.message || 'Failed to fetch articles'}
                    </div>
                )}

                {!isLoading && !error && articles.length === 0 && (
                    <div className="text-center text-gray-500">
                        {searchTerm
                            ? `No articles found for "${searchTerm}"`
                            : 'No articles available'}
                    </div>
                )}

                {!isLoading && !error && articles.length > 0 && (
                    <div className='px-8'>
                        {articles.map(article => (
                            <ArticleCard
                                key={article.id}
                                article={article}
                            />
                        ))}

                        <Pagination
                            currentPage={page}
                            totalPages={totalPages}
                            onPageChange={setPage}
                        />
                    </div>
                )}
            </div>
        </div>
    );
};

export default Home;
