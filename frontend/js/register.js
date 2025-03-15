function register() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
        document.getElementById('error').textContent = '请输入用户名和密码';
        return;
    }

    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.message) {
            alert('注册成功！');
            window.location.href = '/frontend/html/login.html';
        } else {
            document.getElementById('error').textContent = data.error || '注册失败';
        }
    })
    .catch(error => {
        document.getElementById('error').textContent = '服务器错误';
        console.error('Error:', error);
    });
}