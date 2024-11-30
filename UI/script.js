document.getElementById('login-form').addEventListener('submit', function(event) {
    event.preventDefault();
    const userId = document.getElementById('user-id').value;
    if (!userId) {
        showToast("Please enter a valid User ID", "error");
        return;
    }
    const url = `http://localhost:8080/join?id=${encodeURIComponent(userId)}`; 
    fetch(url, {
        method: 'GET',  
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (response.ok) { 
            showToast('Joined successfully!', 'success');
            localStorage.setItem('id',userId)
            setTimeout(() => {
                window.location.href = 'chat.html';
            }, 2000);
        } else {
            showToast('Failed to join. Please try again.', 'error');
        }
    })
    .catch(error => {
        showToast('An error occurred. Please try again later.', 'error');
    });
});


function showToast(message, type) {
    const toast = document.createElement('div');
    toast.classList.add('toast', type);
    toast.textContent = message;

    document.body.appendChild(toast);

    setTimeout(() => {
        toast.remove();
    }, 3000);
}
