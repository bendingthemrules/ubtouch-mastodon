
/*
* Copyright (C) 2022  Development@bendingtherules.nl
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation; version 3.
*
* first is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
import QtQuick 2.7
import Ubuntu.Components 1.3
import Ubuntu.Components.Popups 1.3
import QtQuick.Layouts 1.3
import Qt.labs.settings 1.0

Page {
   id : page
   width : parent.width
   height : parent.height
   anchors.fill : parent
   header : PageHeader {
       id : header
       title : i18n.tr('App Title')
       visible : false
   }
   Rectangle {
       anchors.fill : parent
       color : "#353454"
   }
   Column {
       anchors.left : parent.left
       anchors.right : parent.right
       anchors.top : parent.top
       anchors {
           topMargin : 12
           leftMargin : 20
           rightMargin : 20
       }
       spacing : 20
       Text {
           id : title
           text : "Let's get you setup "+ QClient.selectedServer.domain
           anchors {
               left : parent.left
               right : parent.right
           }
           color : "#F3F2F7"
           font.pointSize : 20
           wrapMode : Text.WordWrap
       }
       Column {
           anchors {
               left : parent.left
               right : parent.right
           }
           spacing : 8
           TextInput {
               id : displayName
               font.pointSize : 12
               text : ""
               color : "#F3F2F7"
               width : parent.width
               padding: 16
               layer.enabled : true
               Label {
                   color : "#9898B2"
                   text : "Display name"
                   font.pointSize : 12
                   anchors.verticalCenter : parent.verticalCenter
                   anchors.left : parent.left
                   anchors.leftMargin : 16
                   visible : !displayName.text && !displayName.activeFocus
               }
               Rectangle {
                   z : -1
                   anchors.fill : parent
                   color : "#2B2937"
                   radius : 8
                   anchors.right : parent.right
                   anchors.left : parent.left
               }
           }
           TextInput {
               id : username
               font.pointSize : 12
               text : ""
               color : "#F3F2F7"
               width : parent.width
               topPadding: 16
               bottomPadding: 16
               leftPadding: 16
               rightPadding: serverName.width + 20
               layer.enabled : true
               Label {
                   color : "#9898B2"
                   text : "Username"
                   font.pointSize : 12
                   anchors.verticalCenter : parent.verticalCenter
                   anchors.left : parent.left
                   anchors.leftMargin : 16
                   visible : !username.text && !username.activeFocus
               }

               Text {
                   id: serverName
                   text: "@" + QClient.selectedServer.domain
                   color: "#F3F2F7"
                   font.weight: Font.Medium
                   anchors.right: parent.right
                   anchors.rightMargin: 16
                   anchors.verticalCenter : parent.verticalCenter
               }

               Rectangle {
                   z : -1
                   anchors.fill : parent
                   color : "#2B2937"
                   radius : 8
                   anchors.right : parent.right
                   anchors.left : parent.left
               }
           }
       }
       Column {
           anchors {
               left : parent.left
               right : parent.right
           }
           spacing : 8
           TextInput {
               id : email
               font.pointSize : 12
               text : ""
               color : "#F3F2F7"
               width : parent.width
               padding: 16
               layer.enabled : true
               Label {
                   color : "#9898B2"
                   text : "Email"
                   font.pointSize : 12
                   anchors.verticalCenter : parent.verticalCenter
                   anchors.left : parent.left
                   anchors.leftMargin : 16
                   visible : !email.text && !email.activeFocus
               }
               Rectangle {
                   z : -1
                   anchors.fill : parent
                   color : "#2B2937"
                   radius : 8
                   anchors.right : parent.right
                   anchors.left : parent.left
               }
           }
           TextInput {
               id : password
               font.pointSize : 12
               text : ""
               color : "#F3F2F7"
               width : parent.width
               padding: 16
               layer.enabled : true
               echoMode: TextInput.Password
               Label {
                   color : "#9898B2"
                   text : "Password"
                   font.pointSize : 12
                   anchors.verticalCenter : parent.verticalCenter
                   anchors.left : parent.left
                   anchors.leftMargin : 16
                   visible : !password.text && !password.activeFocus
               }
               Rectangle {
                   z : -1
                   anchors.fill : parent
                   color : "#2B2937"
                   radius : 8
                   anchors.right : parent.right
                   anchors.left : parent.left
               }
           }
           Text {
               id : passwordDesc
               text : "Include capital letters, special characters, and numbers to increase your password strength."
               anchors {
                   left : parent.left
                   right : parent.right
               }
               color : "#9898B2"
               font.pointSize : 10
               wrapMode : Text.WordWrap
           }
          TextInput {
              id : reason
              visible: QClient.selectedServer.requiresApproval
              font.pointSize : 12
              text : ""
              color : "#F3F2F7"
              width : parent.width
              padding: 16
              layer.enabled : true
              Label {
                  color : "#9898B2"
                  text : "Why do you want to join?"
                  font.pointSize : 12
                  anchors.verticalCenter : parent.verticalCenter
                  anchors.left : parent.left
                  anchors.leftMargin : 16
                  visible : !reason.text && !reason.activeFocus
              }
              Rectangle {
                  z : -1
                  anchors.fill : parent
                  color : "#2B2937"
                  radius : 8
                  anchors.right : parent.right
                  anchors.left : parent.left
              }
          }
       }
   }

   Rectangle {
                id: footer
                height: 64;
                color: "#353454"

                anchors.right: parent.right
                anchors.left: parent.left
                anchors.bottom: parent.bottom
                anchors.rightMargin: 20
                anchors.leftMargin: 20

                z: 2

                RowLayout {
                   anchors.fill: parent

                   Button  {
                       color: "#464766"
                       onClicked: {
                            pageStack.pop()
                       }

                       Label {
                            text: "Back"
                            color: "#F3F2F7"
                            font.weight: Font.Medium
                            anchors.verticalCenter: parent.verticalCenter
                            anchors.horizontalCenter: parent.horizontalCenter
                        }
                    }

                    Item { Layout.fillWidth: true }

                    Button {
                        id: createAccountButton
                        color: "#F3F2F7"
                        onClicked: {
                            const createAccountErrorResponse = QClient.createAccount(displayName.text, username.text, email.text, password.text, reason.text)
                            if(!createAccountErrorResponse){
                                pageStack.clear()
                                pageStack.push(Qt.resolvedUrl("AwaitingActivation.qml"))
                                return
                            }
                            PopupUtils.open(createDialog, root, {'errorResponse': createAccountErrorResponse})
                        }

                        Label {
                            text: "Next"
                            color: "#2B2937"
                            font.weight: Font.Medium
                            anchors.verticalCenter: parent.verticalCenter
                            anchors.horizontalCenter: parent.horizontalCenter
                        }
                    }
                }
    }
   Component {
     id: createDialog
     Dialog {
         id: dialog
         property string errorResponse
         text: errorResponse
         title: "Failed to create"
         Button {
             text: "Continue"
             onClicked: PopupUtils.close(dialog)
         }
     }
   }
}
