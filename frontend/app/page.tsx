"use client";
import { useState, useEffect } from "react";

interface Expression {
  id: string;
  status: string;
  result: string | null;
}

const API_BASE_URL = process.env.ORCHESTRATOR_URL || "http://localhost:8080";

export default function Home() {
  const [expression, setExpression] = useState("");
  const [expressionsDict, setExpressionsDict] = useState<Record<string, Expression>>({});
  const [order, setOrder] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchExpressions();
    const interval = setInterval(fetchExpressions, 5000);
    return () => clearInterval(interval);
  }, []);

  const fetchExpressions = async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/api/v1/expressions`);
      if (!res.ok) {
        handleError(res.status);
        return;
      }
      setError(null);
      const data = await res.json();
      setExpressionsDict((prevDict) => {
        const newDict = { ...prevDict };
        setOrder((prevOrder) => {
          const newOrder = [...prevOrder];
          data.expressions.forEach((expr: Expression) => {
            const exprId = String(expr.id);
            if (!prevOrder.includes(exprId)) {
              newOrder.unshift(exprId);
            }
            newDict[exprId] = expr;
          });
          return newOrder;
        });
        return newDict;
      });
    } catch (error: any) {
      setError(`Ошибка при получении данных: ${error.message}`);
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_BASE_URL}/api/v1/calculate`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ expression }),
      });
      if (!res.ok) {
        handleError(res.status);
        setLoading(false);
        return;
      }
      setExpression("");
      fetchExpressions();
    } catch (error: any) {
      setError(`Ошибка при отправке выражения: ${error.message}`);
    }
    setLoading(false);
  };

  const handleError = (status: number) => {
    let errorMessage = "";
    switch (status) {
      case 400:
        errorMessage = "Неверный запрос. Проверьте корректность.";
        break;
      case 401:
        errorMessage = "Требуется авторизация.";
        break;
      case 404:
        errorMessage = "Ресурс не найден.";
        break;
      case 422:
        errorMessage = "Некорректное выражение, попробуйте другое";
        break;
      case 500:
        errorMessage = "Ошибка сервера. Попробуйте позже.";
        break;
      default:
        errorMessage = "Произошла непредвиденная ошибка.";
    }
    setError(`Ошибка ${status}: ${errorMessage}`);
  };

  const filteredOrder = order.filter((id) => String(id).includes(searchQuery));

  return (
      <div className="min-h-screen bg-gradient-to-r from-blue-100 via-indigo-100 to-purple-100 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900 flex flex-col items-center justify-center p-8">
        <div className="w-full max-w-4xl bg-white dark:bg-gray-800 rounded-lg shadow-lg p-8">
          <h1 className="text-4xl font-bold text-gray-800 dark:text-gray-100 text-center mb-8">
            Распределённый вычислитель
          </h1>
          <form onSubmit={handleSubmit} className="flex flex-col md:flex-row items-center justify-center mb-6">
            <input
                type="text"
                value={expression}
                onChange={(e) => setExpression(e.target.value)}
                className="w-full md:w-2/3 p-3 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-300 dark:bg-gray-700 dark:text-gray-100 mb-4 md:mb-0 md:mr-4"
                placeholder="Введите выражение..."
            />
            <button
                type="submit"
                className="w-full md:w-auto bg-indigo-600 hover:bg-indigo-700 transition-all duration-300 text-white font-semibold py-3 px-6 rounded-lg shadow"
                disabled={loading}
            >
              {loading ? "Отправка..." : "Вычислить"}
            </button>
          </form>
          {error && <div className="mb-6 p-4 bg-red-100 dark:bg-red-900 border border-red-300 dark:border-red-700 text-red-800 dark:text-red-200 rounded">{error}</div>}
          <div className="mb-6">
            <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full p-3 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-300 dark:bg-gray-700 dark:text-gray-100"
                placeholder="Поиск по ID..."
            />
          </div>
          <h2 className="text-2xl font-semibold text-gray-700 dark:text-gray-200 mb-4">Выражения</h2>
          <ul>
            {filteredOrder.map((id) => {
              const expr = expressionsDict[id];
              if (!expr) return null;
              return (
                  <li key={expr.id} className="bg-gray-100 dark:bg-gray-700 rounded-lg shadow p-4 mb-4 transform transition duration-500 hover:scale-105 animate-fadeIn">
                    <p className="text-gray-800 dark:text-gray-200"><strong>ID:</strong> {expr.id}</p>
                    <p className="text-gray-800 dark:text-gray-200"><strong>Статус:</strong> {expr.status}</p>
                    <p className="text-gray-800 dark:text-gray-200"><strong>Результат:</strong> {expr.result || "Ожидайте..."}</p>
                  </li>
              );
            })}
          </ul>
        </div>
        <style jsx>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(20px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .animate-fadeIn { animation: fadeIn 0.5s ease-out; }
      `}</style>
      </div>
  );
}
