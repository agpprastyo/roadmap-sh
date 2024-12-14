import {Link} from 'react-router-dom';
import {RoutePaths} from "../general/RoutePaths.jsx";

const ArticleCard = ({article}) => {
    const formatDate = (dateString) => {
        return new Date(dateString).toLocaleDateString();
    };

    return (
        <div
            className=" shadow-md hover:shadow-lg transition-shadow duration-300 p-4 mb-4 rounded-lg bg-white">
            <Link to={RoutePaths.DETAIL.replace(':id', article.id)}>
                <h2 className="text-xl font-bold mb-2 text-gray-800 hover:text-blue-600 transition-colors duration-200">
                    {article.title}
                </h2>
            </Link>
            <div
                className="mb-2 text-gray-700"
                dangerouslySetInnerHTML={{__html: article.content.substring(0, 200) + '...'}}
            />
            <div className="text-gray-500 text-sm">
                Published on: {formatDate(article.created_at)}
            </div>
        </div>
    );
};

export default ArticleCard;
