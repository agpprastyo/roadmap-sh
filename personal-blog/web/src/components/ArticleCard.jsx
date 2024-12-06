import { Link } from 'react-router-dom';
import {RoutePaths} from "../general/RoutePaths.jsx";


const ArticleCard = ({ article }) => {
    const formatDate = (dateString) => {
        return new Date(dateString).toLocaleDateString();
    };

    return (
        <div className="border p-4 mb-4 rounded">
            <Link to={RoutePaths.DETAIL.replace(':id', article.id)}>
                <h2 className="text-xl font-bold mb-2 hover:text-blue-600">
                    {article.title}
                </h2>
            </Link>
            <div
                className="mb-2"
                dangerouslySetInnerHTML={{ __html: article.content.substring(0, 200) + '...' }}
            />
            <div className="text-gray-500 text-sm">
                Published on: {formatDate(article.created_at)}
            </div>
        </div>
    );
};

export default ArticleCard;
