function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
        document.getElementById('error').textContent = '请输入用户名和密码';
        return;
    }

    fetch('/login', {
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
        if (data.token) {
            localStorage.setItem('token', data.token);
            alert('登录成功！');
            window.location.href = '/game';
        } else {
            document.getElementById('error').textContent = data.error || '登录失败';
        }
    })
    .catch(error => {
        document.getElementById('error').textContent = '服务器错误';
        console.error('Error:', error);
    });
}