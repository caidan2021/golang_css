/*
 * @Date: 2022-09-05 21:01:06
 */
function commonFetch(url, options) {
    return fetch(url, options).then(function(response) {
        if (response.status == 403) {
            location.href = "/admin/view/admin_login"
        }
        return response.json();
    // }).catch(function(error) {
    //     console.log(error.message)
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

function renderSideBar()
{
    // url = "/admin/menu";
    // options = {
    //     method: 'get',
    //     headers: {
    //         "Content-Type": "application/json charset=utf-8"
    //     },
    // }
    // var menuTpl = "<li class='nav-item '><a class='nav-link active' href='#' data-toggle='collapse' aria-expanded='false' data-target='#submenu-{parentId}' aria-controls='submenu-{parentId}'><i class='{icon}'></i>{name} <span class='badge badge-success'>6</span></a>"
    // menuTpl += "<div id='submenu-{parentId}' class='collapse submenu' style=''><ul class='nav flex-column'>{subMenu}</ul></div></li>"
    // subMenuTpl = "<li class='nav-item'><a class='nav-link' href='{pagePath}'>{name}</a></li>"
    // subMenu = "";

    // commonFetch(url, options).then(function (jsonResult) {
    //     list = jsonResult.data.items;
    //     // 先处理树
    //     // list = getTree(list)
    //     $.each(list, function (index, item) {
    //         subMenu += subMenuTpl.replace(/{(\w+)}/g, function(_, $1) {
    //             return item[$1]
    //         })
    //     })
    //     console.log(menu)
    //     console.log(subMenu)
    //     menu = menu.replace("{subMenu}", subMenu)
    //     console.log(menu)
    //     $(".sidebar").append(menu)
    // })

}