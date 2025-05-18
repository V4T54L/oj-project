import { useState } from "react";
import { Modal, ModalHeader, ModalBody, ModalFooter } from "./Modal";

export const UserFormModal = () => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <>
            <button
                onClick={() => setIsOpen(true)}
                className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 transition text-white rounded"
            >
                Open Modal
            </button>

            <Modal
                isOpen={isOpen}
                onOpenChange={setIsOpen}
                size="md"
                radius="lg"
                backdrop="blur"
                placement="center"
                scrollBehavior="inside"
            >
                <ModalHeader>
                    Create New Account
                </ModalHeader>

                <ModalBody>
                    <form className="space-y-4">
                        <div>
                            <label className="block text-sm text-gray-300 mb-1">Name</label>
                            <input
                                type="text"
                                className="w-full bg-gray-700 border border-gray-600 text-white rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm text-gray-300 mb-1">Email</label>
                            <input
                                type="email"
                                className="w-full bg-gray-700 border border-gray-600 text-white rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                                required
                            />
                        </div>
                    </form>
                </ModalBody>

                <ModalFooter>
                    <button
                        onClick={() => setIsOpen(false)}
                        className="text-gray-400 hover:text-gray-200 transition"
                    >
                        Cancel
                    </button>
                    <button className="bg-indigo-600 hover:bg-indigo-700 transition text-white px-4 py-2 rounded">
                        Submit
                    </button>
                </ModalFooter>
            </Modal>
        </>
    );
};
