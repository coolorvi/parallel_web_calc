document.addEventListener("DOMContentLoaded", function () {
    const expressionForm = document.getElementById("expression-form");
    const expressionInput = document.getElementById("expression-input");
    const resultContainer = document.getElementById("result");
    const tasksContainer = document.getElementById("tasks");
    
    expressionForm.addEventListener("submit", async function (event) {
        event.preventDefault();
        const expression = expressionInput.value.trim();
        if (!expression) {
            alert("Введите выражение");
            return;
        }
        
        try {
            const response = await fetch("http://localhost:8080/evaluate", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ expression })
            });
            
            if (!response.ok) {
                throw new Error("Ошибка при отправке запроса");
            }
            
            const data = await response.json();
            resultContainer.innerText = `ID выражения: ${data.expression_id}`;
        } catch (error) {
            console.error("Ошибка:", error);
            alert("Не удалось отправить выражение");
        }
    });
    
    async function fetchResults() {
        try {
            const response = await fetch("http://localhost:8080/results");
            if (!response.ok) {
                throw new Error("Ошибка при получении результатов");
            }
            
            const data = await response.json();
            tasksContainer.innerHTML = "";
            
            data.results.forEach(result => {
                const div = document.createElement("div");
                div.className = "task-result";
                div.innerText = `Выражение ID: ${result.expression_id}, Результат: ${result.result}`;
                tasksContainer.appendChild(div);
            });
        } catch (error) {
            console.error("Ошибка:", error);
        }
    }
    
    setInterval(fetchResults, 5000);
});
