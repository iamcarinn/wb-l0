// script.js
document.getElementById('orderForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent form from submitting the traditional way
    
    const orderUID = document.getElementById('orderUID').value;

    // Отправляем POST-запрос на сервер
    fetch('http://127.0.0.1:3330/orders', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ order_uid: orderUID }),
    })
    .then(response => response.json())  // Преобразуем ответ в формат JSON
    .then(data => {
        document.getElementById('result').textContent = JSON.stringify(data, null, 2);
    })
    .catch(error => {
        document.getElementById('result').textContent = 'Error: ' + error.message;
    });
});
