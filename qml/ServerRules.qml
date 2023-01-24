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
import QtQuick.Layouts 1.3
import Qt.labs.settings 1.0

Page {
    id: page
    width: parent.width
    height: parent.height
    anchors.fill: parent

    header: PageHeader {
        id: header
        title: i18n.tr('App Title')
        visible: false
    }

    Rectangle {
        anchors.fill: parent
        color: "#353454"
    }

    ListView {
        id: mastodonRules

        anchors.fill: parent
        anchors {
            topMargin: 12
            leftMargin: 20
            rightMargin: 20
        }

        header: Column {
            id: rulesHeader
            width: parent.width
            spacing: 12

            Text {
                id: title
                text: "Some ground rules"
                anchors {
                     left: parent.left
                     right: parent.right
                }

                color: "#F3F2F7"
                font.pointSize: 20
                wrapMode: Text.WordWrap
            }

            Text {
                id: info
                text: "Take a minute to review the rules set and enforced by " + QClient.selectedServer.domain + "."
                anchors {
                     left: parent.left
                     right: parent.right
                }

                wrapMode: Text.WordWrap
                font.pointSize: 12
                font.weight: Font.Medium
                color: "#9898B2"
            }

            Item {height: 4; width: 1} // spacer
        }

        model: QClient.selectedServer.serverRules
        delegate: ColumnLayout {
            spacing: 16

            RowLayout {
                id: serverRow
                spacing: 12

                anchors {
                    topMargin: 12
                    bottomMargin: 12
                }

                Rectangle {
                    width: 20
                    height: 20
                    color: "#9898B2"
                    border.width: 2
                    border.color: "#9898B2"
                    radius: 100

                    Text {
                        text: index
                        color: "#353454"
                        font.weight: Font.Medium
                        anchors.verticalCenter: parent.verticalCenter
                        anchors.horizontalCenter: parent.horizontalCenter
                    }
                }

                Text {
                    Layout.preferredWidth: mastodonRules.width - parent.spacing - 20 // icon size
                    wrapMode: Text.WordWrap
                    text: QClient.selectedServer.serverRules[index]
                    color: "#F3F2F7"
                    font.weight: Font.Medium
                }
        }

        Rectangle {
            height: 1
            width: parent.width
            Layout.preferredWidth: mastodonRules.width
            color: "#9898B2"
            visible: index < mastodonRules.count -1
        }

        Item { height: 0; width: 1 } // spacer
    }

        footer: Rectangle {
            id: footer
            width: parent.width;
            height: 64;

            color: "#353454"
            z: 2

            RowLayout {
                anchors.fill: parent

                Button  {
                    color: "#464766"
                    Layout.preferredHeight: backLabel.paintedHeight + 24
                    onClicked: {
                        pageStack.pop()
                    }

                    Label {
                        id: backLabel
                        text: "Back"
                        color: "#F3F2F7"
                        font.weight: Font.Medium
                        anchors.verticalCenter: parent.verticalCenter
                        anchors.horizontalCenter: parent.horizontalCenter
                    }
                }

                Item { Layout.fillWidth: true }

                Button {
                    color: "#F3F2F7"
                    Layout.preferredHeight: nextLabel.paintedHeight + 24
                    onClicked: {
                        if(QClient.getIsRegistering()){
                            pageStack.push(Qt.resolvedUrl("CreateAccount.qml"))
                        }
                        if(QClient.getIsLoggingIn()){
                            QClient.setLoginUrl()
                            pageStack.push(Qt.resolvedUrl("Webview.qml"))
                        }
                    }

                    Label {
                        id: nextLabel
                        text: "Next"
                        color: "#2B2937"
                        font.weight: Font.Medium
                        anchors.verticalCenter: parent.verticalCenter
                        anchors.horizontalCenter: parent.horizontalCenter
                    }
                }
            }
        }
        footerPositioning: ListView.OverlayFooter
    }
}