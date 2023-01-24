
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
import Ubuntu.Components 1.3 // import QtQuick.Controls 2.2
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
        anchors.centerIn: parent
        width: parent.width - 40
        spacing: 16
          Text {
              id : title
              text : "One last thing"
              anchors {
                  left : parent.left
                  right : parent.right
              }
              color : "#F3F2F7"
              font.pointSize : 20
              wrapMode : Text.WordWrap
          }

          Text {
              id: mailInfo
              text: "Tap the link we emailed to you to verify your account."
              color: "#9898B2"
              font.pointSize : 12
              wrapMode : Text.WordWrap
              anchors {
                  left : parent.left
                  right : parent.right
              }
          }

          Image {
              id: image

              width: parent.width
              height: parent.width
              Layout.alignment: Qt.AlignLeft | Qt.AlignTop
              Layout.fillHeight: false
              Layout.fillWidth: true
              source: "qrc:/assets/email.png"
              fillMode: Image.PreserveAspectFit
          }
      }
   Timer {
       interval: 1000
       running: true
       repeat: true
       onTriggered: () => {
            if(QClient.awaitingActivation){
                return
            }
            pageStack.clear()
            pageStack.push(Qt.resolvedUrl("Webview.qml"))
       }
   }
}