// Функция проверки здоровья бэкенда
async function checkHealth() {
    try {
        const res = await fetch('/health');
        const data = await res.json();
        const statusEl = document.getElementById('health-status');
        statusEl.textContent = `Статус: ${data.status} | БД: ${data.db ? 'OK' : 'Ошибка'}`;
        statusEl.style.color = data.db ? 'green' : 'red';
    } catch (e) {
        console.error('Ошибка проверки health:', e);
    }
}

// Функция загрузки списка заметок
async function fetchNotes() {
    try {
        const res = await fetch('/api/notes');
        const notes = await res.json();
        const list = document.getElementById('notes-list');
        
        if (notes.length === 0) {
            list.innerHTML = '<li>Заметок пока нет.</li>';
            return;
        }

        list.innerHTML = notes.map(n => 
            `<li><strong>${n.title}</strong>: ${n.body}</li>`
        ).join('');
    } catch (e) {
        console.error('Ошибка загрузки заметок:', e);
    }
}

// Обработка отправки формы
document.getElementById('note-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const title = document.getElementById('title').value;
    const body = document.getElementById('body').value;

    try {
        await fetch('/api/notes', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title, body })
        });

        // Очищаем форму и обновляем список
        document.getElementById('title').value = '';
        document.getElementById('body').value = '';
        fetchNotes();
    } catch (e) {
        alert('Не удалось добавить заметку');
        console.error(e);
    }
});

// Запускаем при загрузке страницы
checkHealth();
fetchNotes();
