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
//import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
import Qt.labs.settings 1.0

Page {
    id: page
    width: parent.width
    height: parent.height
    anchors.fill: parent
    implicitHeight: 0;

    header: PageHeader {
        id: header
        title: i18n.tr('App Title')
        visible: false
    }

    Image {
        id: background
        fillMode: Image.PreserveAspectCrop
        source: "qrc:/assets/mastodonBG.png"
        antialiasing: true
        anchors.fill: parent
    }


    Image {
        id: image
        fillMode: Image.PreserveAspectFit
        source: "qrc:/assets/mastodonLogo.svg"

        anchors {
            top: header.top
            left: parent.left
            right: parent.right

            topMargin: 48
            leftMargin: 40
            rightMargin: 40
        }

        verticalAlignment: Label.AlignVCenter
        horizontalAlignment: Label.AlignHCenter
    }

    Label {
        font.weight: Font.Bold
        anchors {
            top: header.top
            left: parent.left
            right: parent.right
        }

        verticalAlignment: Label.AlignVCenter
        horizontalAlignment: Label.AlignHCenter
    }

    ColumnLayout {
            id: columnLayout
            spacing: 24

            anchors {
                left: parent.left
                right: parent.right
                bottom: parent.bottom

                leftMargin: 40
                rightMargin: 40
                bottomMargin: 48
            }

            Button {
                id: getStarted
                Layout.fillWidth: true
                Layout.maximumWidth: 600
                Layout.preferredHeight: startedLabel.paintedHeight + 24
                Layout.alignment: Qt.AlignHCenter | Qt.AlignVCenter

                color: "#F3F2F7"
                onClicked: {
                    QClient.setIsRegistering()
                    pageStack.push(Qt.resolvedUrl("Servers.qml"))
                }
                Label {
                    id: startedLabel
                    text: "Get started"
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.horizontalCenter: parent.horizontalCenter
                    fontSize: "medium"
                    font.weight: Font.Medium
                    color: "#2B2937"
                }
            }

            Button {
                id: login
                Layout.fillWidth: true
                Layout.maximumWidth: 600
                Layout.preferredHeight: loginLabel.paintedHeight + 24
                Layout.alignment: Qt.AlignHCenter | Qt.AlignVCenter

                color: "#57AC82"

                onClicked: {
                    QClient.setIsLoggingIn()
                    pageStack.push(Qt.resolvedUrl("Servers.qml"))
                }
                Label {
                    id: loginLabel
                    text: "Log in"
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.horizontalCenter: parent.horizontalCenter
                    fontSize: "medium"
                    font.weight: Font.Medium
                    color: "#ffffff"
                }
            }
        }
}