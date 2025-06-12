import React from 'react';

const TestCaseEditor: React.FC<any> = ({ testCases, activeTab, setTestCases, setActiveTab }) => (
    <div>
        <h2 className="font-semibold mb-2">Test Cases</h2>
        <div className="flex gap-2 mb-2">
            {testCases?.map((_, idx) => (
                <button
                    key={idx}
                    onClick={() => setActiveTab(idx)}
                    className={`px-3 py-1 rounded border ${activeTab === idx ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 border-gray-300'}`}
                >
                    Case {idx + 1}
                </button>
            ))}
        </div>
        <div className="space-y-2">
            {
                testCases?.length > 0 &&
                <>
                    <textarea
                        className="w-full p-2 border rounded font-mono text-sm"
                        rows={3}
                        value={testCases[activeTab]?.Input || ''}
                        onChange={(e) => {
                            const updated = [...testCases];
                            updated[activeTab].Input = e.target.value;
                            setTestCases(updated);
                        }}
                        placeholder="Input"
                    />
                    <textarea
                        className="w-full p-2 border rounded font-mono text-sm"
                        rows={3}
                        value={testCases[activeTab]?.ExpectedOutput || ''}
                        onChange={(e) => {
                            const updated = [...testCases];
                            updated[activeTab].ExpectedOutput = e.target.value;
                            setTestCases(updated);
                        }}
                        placeholder="Expected Output"
                    />
                </>
            }
        </div>
    </div>
);

export default TestCaseEditor;
