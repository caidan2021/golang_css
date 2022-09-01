
function adminRegister()
{
    url = "/admin/user/add"
    options = {
        method: 'post',
        headers: {
            "Content-Type": "application/json charset=utf-8"
        },
        body: JSON.stringify({
            'name': $('#name').val(),
            'password': $('#pass').val()
        })
    }

    fetch(url, options).then(function(response) {
        return response.json();
    }).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("添加用户失败, err: " + jsonResult.msg);
            return
        }
        location.href = "/admin/view/admin_login.html"
    }).catch(function(error) {
        console.log(error);
        alert("添加用户失败, err: " + error.message)
    });
}

function adminLogion()
{
    url = "/admin/user/login";
    options = {
        method: 'post',
        headers: {
            "Content-Type": "application/json charset=utf-8"
        },
        body: JSON.stringify({
            'username': $('#username').val(),
            'password': $('#password').val(),
            'rememberMe': $('#remember').val() == 'on' ? true : false
        })
    }
    fetch(url, options).then(function(response) {
        return response.json();
    }).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("登陆失败, err: " + jsonResult.msg)
            return
        }
        location.href = "/admin/view/admin_dashboard.html"
    }).catch(function(error) {
        console.log(error);
        alert("登陆失败, err: " + error.message);
    });
}

