/*
 * @Date: 2022-08-31 11:10:13
 */

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
        alert("登陆失败, err: " + error.message);
    });
}

function currentUser()
{
    url = "/admin/user/current";
    options = {
        method: 'get',
        headers: {
            "Content-Type": "application/json charset=utf-8"
        },
    }
    commonFetch(url, options).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("获取用户失败, err: " + jsonResult.msg)
            return
        }
        $("#currentUserName").html(jsonResult.data.currentUser.name)
        $("#currentUserAvatar").attr("src", jsonResult.data.currentUser.avatar)
        return
    }).catch(function(error) {
        alert("获取用户失败, err: " + error.message);
    });

}

