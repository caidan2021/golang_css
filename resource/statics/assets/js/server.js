/*
 * @Date: 2022-08-23 14:45:04
 */

function commonFetch(url, options) {
    return fetch(url, options).then(function(response) {
        if (response.status == 403) {
            location.href = "/admin/view/admin_login"
            return
        }
        return response.json();
    })
}

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


function orderList() {
    url = "/admin/order/list";

    request = GetRequest()
    if (request.orderNo != undefined) {
        url = url + "?orderNo=" + request.orderNo
    }

    options = {
        method: 'get',
        headers: {
            "Content-Type": "application/json charset=utf-8",
            "Css-Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InRlc3QiLCJleHAiOjE2NjI0NzQ5MTR9.hu6Uem_En5QWReoy81-bioQBOZ_BMFlAUpvjImqSk34"
        },
    }

    commonFetch(url, options).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("获取订单失败, err: " + jsonResult.msg)
            return
        }

        var html = "<thead><tr><th width='5%'>ID</th><th>图片</th><th width='15%'>订单No</th><th width='10%'>创建时间</th></tr></thead><tbody id='tbody'></tbody>"
        var tr = "<tr><td>{id}</td><td>{thumbnail}</td><td>{outOrderNo}</td><td>{createdTime}</td></tr>"

        list = jsonResult.data.list
        $.each(list, function(index, item) {
            html += tr.replace(/{(\w+)}/g, function(_, $1) {
                if ($1 == 'thumbnail') {
                    var thumbnails = "";
                    for (k = 0; k < item[$1].length; k++) {
                        if (k != 0) {
                            thumbnails += "&&<img src='" + item[$1][k] + "' height='75px' alt='加载失败' />"
                        } else {
                            thumbnails += "<img src='" + item[$1][k] + "' height='75px' alt='加载失败' />"
                        }
                    }
                    return thumbnails
                }
                return item[$1]
            });

        });
        $(".order-list").append(html)
    }).catch(function(error) {
        console.log(error);
    });

}


