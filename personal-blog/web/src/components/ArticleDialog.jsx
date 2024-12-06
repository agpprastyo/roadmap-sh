import { useState, useEffect } from "react";

const ArticleDialog = ({
                           isOpen,
                           onClose,
                           onSubmit,
                           initialData = null,
                           dialogMode = 'create'
                       }) => {
    // Initial state
    const INITIAL_STATE = {
        title: '',
        content: '',
        published: false
    };

    const [formData, setFormData] = useState(INITIAL_STATE);

    // Update form data when initialData changes (for edit mode)
    useEffect(() => {
        if (initialData) {
            setFormData({
                title: initialData.title || '',
                content: initialData.content || '',
                published: initialData.published || false
            });
        }
    }, [initialData]);

    // Reset form when dialog is closed or opened
    useEffect(() => {
        if (!isOpen) {
            // Reset to initial state when dialog is closed
            setFormData(initialData ?
                {
                    title: initialData.title || '',
                    content: initialData.content || '',
                    published: initialData.published || false
                } :
                INITIAL_STATE
            );
        }
    }, [isOpen, initialData]);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prevState => ({
            ...prevState,
            [name]: type === 'checkbox' ? checked : value
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        // Call onSubmit with form data and article ID (for edit mode)
        onSubmit(initialData ? initialData.id : null, {
            title: formData.title,
            content: formData.content,
            published: formData.published
        });

        // Reset form after submission
        setFormData(INITIAL_STATE);
    };

    const handleCancel = () => {
        // Reset form
        setFormData(INITIAL_STATE);
        // Close the dialog
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
            <div className="bg-white p-6 rounded-lg w-96">
                <h2 className="text-xl font-bold mb-4">
                    {dialogMode === 'create' ? 'Create New Article' : 'Edit Article'}
                </h2>
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label className="block mb-2">Title</label>
                        <input
                            type="text"
                            name="title"
                            value={formData.title}
                            onChange={handleChange}
                            className="w-full border p-2 rounded"
                            required
                        />
                    </div>
                    <div className="mb-4">
                        <label className="block mb-2">Content</label>
                        <textarea
                            name="content"
                            value={formData.content}
                            onChange={handleChange}
                            className="w-full border p-2 rounded"
                            rows="4"
                            required
                        />
                    </div>
                    <div className="mb-4">
                        <label className="flex items-center">
                            <input
                                type="checkbox"
                                name="published"
                                checked={formData.published}
                                onChange={handleChange}
                                className="mr-2"
                            />
                            Publish Article
                        </label>
                    </div>
                    <div className="flex justify-end space-x-2">
                        <button
                            type="button"
                            onClick={handleCancel}
                            className="bg-gray-300 text-black p-2 rounded"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="bg-blue-500 text-white p-2 rounded"
                        >
                            {dialogMode === 'create' ? 'Create Article' : 'Update Article'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default ArticleDialog;
