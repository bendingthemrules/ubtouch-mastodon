import QtQuick 2.9
import QtQuick.Layouts 1.1
import Ubuntu.Components 1.3
import Ubuntu.Content 1.3
import Ubuntu.Components.Popups 1.3
import Qt.labs.settings 1.0
import QtWebEngine 1.10
import "../qml"

Page {
    id: webview
    header: PageHeader {
       id: header
       title: i18n.tr('App Title')
       visible: false
    }

    property QtObject defaultProfile: WebEngineProfile {
        id: webContext
        httpUserAgent: "Mozilla/5.0 (Linux; Android 11; Ubuntu Touch) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.72 Mobile Safari/537.36"
        downloadPath: QClient.appDir + "Downloads/"

        onDownloadRequested: {
            console.log("a download was requested with path %1".arg(download.path))
            console.log(download.url.toString())
            download.accept();
        }

        onDownloadFinished: {
            console.log("a download was finished with path %1.".arg(download.path))
        }
    }

    Component {
        id: filePickerComponent
        PickerDialog {}
    }
    Component {
        id: changeServersDialog
        Dialog {
            id: dialog
            text: "You're about to change servers, do you want to continue?"
            title: "Change servers"
            Button {
                text: "Continue"
                color: "#595aff"
                onClicked: () => {
                    PopupUtils.close(dialog)
                    pageStack.push(Qt.resolvedUrl("StartScreen.qml"))
                }
            }
            Button {
                text: "Cancel"
                onClicked: () => {
                    PopupUtils.close(dialog)
                }
            }
        }
    }

    WebEngineView {
        id: webEngineView
        width: parent.width
        height: parent.height
        visible: false

        profile: defaultProfile

        settings.javascriptCanAccessClipboard: true

        onLoadProgressChanged: {
            progressBar.value = loadProgress
            if (loadProgress === 100){
                changeServers.visible = false;
                console.log('Url changed to '+ url)
                if(url.toString().includes("/auth/sign_in") || url.toString().includes("/about")) {
                    console.log('Setting change servers button to visible')
                    changeServers.visible = true
                }
                webEngineView.runJavaScript("insertChangeServers = function(){
                                                 if(document.querySelector('#internal-change-servers')){
                                                     return;
                                                 }
                                                 const node = document.createElement('a');
                                                 node.classList='column-link column-link--transparent';
                                                 node.id = 'internal-change-servers'
                                                 node.title='change-server';
                                                 node.href='/internal-change-servers';
                                                 const innerNode = document.createElement('i');
                                                 innerNode.classList='fa fa-server column-link__icon fa-fw';
                                                 node.appendChild(innerNode);
                                                 const spanNode = document.createElement('span');
                                                 spanNode.innerText = 'Change server';
                                                 node.appendChild(spanNode);
                                                 const target = document.querySelector('.navigation-panel .flex-spacer');
                                                 if(!target){
                                                     return;
                                                 }
                                                 target.parentNode.insertBefore(node, target); node.onclick;

                                             }; insertChangeServers() ")
                QClient.setProfile(webEngineView.profile)
                visible = true;
            }
        }
        onLinkHovered: (hoveredUrl) => {console.log(hoveredUrl)}
        zoomFactor: 1
        anchors.fill: parent
        url: QClient.webviewUrl
        onFileDialogRequested: {
            console.log('Requested file dialog')
			request.accepted = true;
			var fakeModel = {
					allowMultipleFiles: request.mode == FileDialogRequest.FileModeOpenMultiple,
					reject: function() {
						request.dialogReject();
					},
					accept: function(files) {
							request.dialogAccept(files);
					}
			};
			var  pickerInstance = filePickerComponent.createObject(webEngineView,{model:fakeModel});
	    }

        // Open external URL's in the browser and not in the app
        onNavigationRequested: (navigationRequest) => {
            console.log('Navigation requested ' + navigationRequest.url)
            if(navigationRequest.url.toString().includes("/internal-change-servers")){
                navigationRequest.action = WebEngineNavigationRequest.IgnoreRequest
                PopupUtils.open(changeServersDialog, root)
            }
        }
        onNewViewRequested: (request) => {
            Qt.openUrlExternally(request.requestedUrl);
        }
    }

    Button  {
        id: changeServers
        color: "#595aff"
        visible: false
        height: changeServersLabel.paintedHeight + 24
        width: changeServersLabel.paintedWidth + 24
        anchors.top: parent.top
        anchors.right: parent.right
        anchors.topMargin: 10
        anchors.rightMargin: 10
        onClicked: {
            pageStack.push(Qt.resolvedUrl("StartScreen.qml"))
        }
        Label {
            id: changeServersLabel
            text: "Change server"
            color: "#F3F2F7"
            font.weight: Font.Medium
            anchors.verticalCenter: parent.verticalCenter
            anchors.horizontalCenter: parent.horizontalCenter
        }
    }

    Rectangle {
        visible: !webEngineView.visible
        color: "#000000"
        anchors.fill: parent
    }
    Column {
        anchors.fill: parent
        visible: !webEngineView.visible
           Image {
              id: image
              anchors.centerIn: parent
              width: sourceSize.width
              height: sourceSize.height
              Layout.alignment: Qt.AlignLeft | Qt.AlignTop
              source: "qrc:/assets/logo.svg"
           }
            ProgressBar {
                id: progressBar
                value: 0
                minimumValue: 0
                maximumValue: 100
                anchors.top: image.bottom
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.topMargin: 30
            }
      }
}
