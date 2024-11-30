

document.addEventListener("DOMContentLoaded", function () {
    const chatBox = document.getElementById('chat-box');
    const sendBtn = document.getElementById('send-btn');
    const messageInput = document.getElementById('message');
    const userId = localStorage.getItem('id');
    if (!userId) {
        alert("User ID not found. Please join the room first.");
        window.location.href = 'index.html'; 
    }

    fetchMessages();

    sendBtn.addEventListener('click', function () {
        const message = messageInput.value.trim();
        if (message) {
            sendMessage(message);
            messageInput.value = ''; 
        }
    });
    function sendMessage(message) {
        const userId = localStorage.getItem('id'); 
        if (!userId) {
            showToast('User ID not found. Please join the room first.', 'error');
            return;
        }

        const url = `http://localhost:8080/send?id=${encodeURIComponent(userId)}&message=${encodeURIComponent(message)}`;
    
        fetch(url, {
            method: 'GET', 
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => {
            if (response.ok) { 
                showToast('Message sent successfully!', 'success');
                fetchMessages(); 
            } else {
                showToast('Failed to send message. Try again later.', 'error');
            }
        })
        .catch(error => {
            console.error('Error sending message:', error);
            showToast(error.message, 'error');
        });
    }
    

    function fetchMessages() {
        fetch(`http://localhost:8080/messages?id=${userId}`)
            .then(response => response.json())
            .then(data => {
                if (data.message != '' ) {
                    console.log('No new messages.')
                }
                displayMessages(data.messages);
            })
            .catch(error => {
                console.error('Error fetching messages:', error);
            });
    }

    function displayMessages(messages) {
        const chatBox = document.getElementById('chat-box');
        const userId = localStorage.getItem('id');
    
        chatBox.innerHTML = ''; 
    
        messages.forEach(msg => {
            const messageContainer = document.createElement('div');
            messageContainer.classList.add('message-container');
    
            if (msg.id === userId) {
                console.log(msg.id)
                messageContainer.classList.add('self');
            } else {
                messageContainer.classList.add('other');
            }
    
            const userIdElement = document.createElement('span');
            userIdElement.classList.add('user-id');
            userIdElement.textContent = `User: ${msg.id}`;
    
            const messageElement = document.createElement('div');
            messageElement.classList.add('message');
            messageElement.textContent = msg.message;
    
            messageContainer.appendChild(userIdElement);
            messageContainer.appendChild(messageElement);
            chatBox.appendChild(messageContainer);
        });
        chatBox.scrollTop = chatBox.scrollHeight;
    }
    
    function showToast(message, type) {
        const toast = document.createElement('div');
        toast.classList.add('toast', type);
        toast.textContent = message;
        document.body.appendChild(toast);

        setTimeout(() => {
            toast.remove();
        }, 3000);
    }
    setInterval(fetchMessages, 1000); 
});










document.getElementById('leave-chat').addEventListener('click', function () {
    const confirmLeave = confirm('Are you sure you want to leave the chat?');
    if (!confirmLeave) {
        return; 
    }

    const userId = localStorage.getItem('id'); 
    const url = `http://localhost:8080/leave?id=${encodeURIComponent(userId)}`;

    fetch(url, {
        method: 'GET',
    })
    .then(response => {
        if (response.ok) {
            console.log('hi'+response.ok)
            showToast('You have left the chat.', 'success');
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 20);
        } else {
            showToast('Failed to leave chat. Please try again.', 'error');
        }
    })
    .catch(error => {
        console.error('Error leaving chat:', error);
        showToast('An error occurred. Please try again.', 'error');
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

