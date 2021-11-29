console.log("+++ index.js +++ 2021-11-29 16:12");

let file2upload = null;
let itemListView = null;

const btnUpload = document.getElementById("btn_upload");
if (btnUpload != null) {
    btnUpload.onclick = () => {
        console.log("Click button to upload file");
        uploadFileToServer(file2upload);
    }
}

function uploadFileToServer(file) {
    console.log(">> uploadFileToServer");
    let formData = new FormData();
    formData.append("file2upload", file)

    fetch('//localhost:8080/files/' + file.name, { // Your POST endpoint
        method: 'POST',
        body: formData // This is your file object
    }
    ).then(
        success => {
            console.log("OK>", success)
            getFilesOnServer();
        } // Handle the success response object
    ).catch(
        error => console.log("FAIL>", error) // Handle the error response object
    );
}

function getFilesOnServer() {
    console.log(">> getFilesOnServer");
    fetch('//localhost:8080/files', {
        method: 'GET'
    }).then(respnose => respnose.json())
        .then(respnose => {
            console.log('Response', respnose)
            showListItems(respnose.files)
        }).catch(error => console.error('Error', error))
}

function deleteFileOnServer(fileName) {
    console.log(">> deleteFileOnServer", fileName);
    fetch('//localhost:8080/files/' + fileName, {
        method: 'DELETE'
    })
        .then(respnose => {
            console.log('Response', respnose)
            getFilesOnServer();
        }).catch(error => console.error('Error', error))
}

function initListItem() {
    let listItems = document.querySelectorAll(".list-group-item");
    itemListView = document.getElementById("item_list_view");
    console.log(">> initListItem", listItems.length);
    if (listItems.length > 0) {
        defaultOfListItemView = listItems[0].cloneNode(true);

        for (let i = 0; i < listItems.length; i++) {
            console.log("Remove:");
            itemListView.removeChild(listItems[i]);
        }
    }
}

function showListItems(files) {
    console.log(">> showListItems", files);

    initListItem();

    let button2Delete = document.getElementById("button_delete_template");

    for (let i = 0; i < files.length; i++) {
        let fileName = files[i];
        let listItem = defaultOfListItemView.cloneNode(true);
        let button = button2Delete.cloneNode(true);
        button.onclick = deleteFileOnServer.bind(this, fileName);
        button.style.visibility = "visible";
        button.innerText = "Delete";
        console.log("List file:", fileName);
        listItem.innerText = fileName;
        listItem.style.visibility = "visible";
        listItem.appendChild(button);
        itemListView.appendChild(listItem);
    }
}

const fileUploader = document.querySelector('#file-uploader');
fileUploader.addEventListener('change', (e) => {
    console.log("fileUploader.onChanged:", e.target.files); // get file object
    file2upload = e.target.files[0];
});

initListItem();
getFilesOnServer();