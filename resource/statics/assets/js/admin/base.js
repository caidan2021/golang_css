/*
 * @Date: 2022-09-05 21:01:06
 */

document.write('<script src="/statics/assets/vendor/jquery/jquery-3.3.1.min.js"></script>')

//popper.js必须在bootstrap.js之前引用进来，否则报错的
document.write('<script src="/statics/assets/vendor/bootstrap/js/popper.min.js"></script>')
document.write('<script src="/statics/assets/vendor/bootstrap/js/bootstrap.bundle.min.js"></script>')
document.write('<script src="/statics/assets/vendor/bootstrap/js/bootstrap.min.js"></script>')
document.write('<script src="/statics/assets/vendor/bootstrap/js/fileinput.min.js"></script>')

document.write('<script src="/statics/assets/vendor/slimscroll/jquery.slimscroll.js"></script>')
document.write('<script src="/statics/assets/vendor/multi-select/js/jquery.multi-select.js"></script>')
document.write('<script src="/statics/assets/libs/js/main-js.js"></script>')

document.write('<script type="text/javascript" src="/statics/assets/js/helper.js"></script>')
document.write('<script type="text/javascript" src="/statics/assets/js/admin/order.js"></script>')
document.write('<script type="text/javascript" src="/statics/assets/js/admin/admin_user.js"></script>')

function commonFetch(url, options) {
    return fetch(url, options).then(function(response) {
        if (response.status == 403) {
            location.href = "/admin/view/admin_login"
            return
        }
        return response.json();
    })
}

function listPaginate(totalCount, pageSize = 15)
{

    paginate = "<nav aria-label='Page navigation'><ul class='pagination justify-content-center'>";
    pageLimit = totalCount % pageSize == 0 ? parseInt(totalCount / pageSize) : parseInt(totalCount / pageSize) + 1;

    currentPage = GetUrlParam("page") ?? 1
    for (i = 0; i < pageLimit; i++) {

        params = [{"name": "page","value": i + 1}]
        cls = 'page-link ' + (currentPage == i + 1 ? 'bg-primary' : '');
        paginate += "<li class='page-item'><a class='" + cls + "' href='" + ReplaceOrAddRequest(params) +  "'>" + (i + 1) + "</a></li>"
    }
    paginate += "</ul></nav>";
    $("#pagination").append(paginate)
}