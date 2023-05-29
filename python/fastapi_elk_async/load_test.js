import http from 'k6/http';


function debug_obj(obj) {
    return JSON.stringify(obj, null, 2)
}

function async() {
    let res = http.get("http://localhost:8000/");
    if (res.status != 200) {
        throw new Error('Error  ' + debug_obj(res))
    }

}

function notasync() {
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    let res = http.get("http://localhost:8000/notasync", params);
    if (res.status != 200) {
        throw new Error('Error  ' + debug_obj(res))
    }
}

export default function () {
    notasync();
    // async();
}
