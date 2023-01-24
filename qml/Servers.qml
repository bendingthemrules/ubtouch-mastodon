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

    VisualDataModel {
        id: visualModel
        model: QClient

        groups: [
            VisualDataGroup {
                name: "matchingSearchTerm"
                includeByDefault: true
            }
        ]
        filterOnGroup: "matchingSearchTerm"
        delegate: ColumnLayout {
                              spacing: 16

                              RowLayout {
                                  id: serverRow
                                  spacing: 12

                                  anchors {
                                      topMargin: 12
                                      bottomMargin: 12
                                  }

                                  CheckBox {
                                      id: mastodonServerCheckbox
                                      checked: selected
                                      width: 0
                                      height: 0

                                      Rectangle {
                                          anchors.fill: parent
                                          color: "#353454"
                                      }

                                      Rectangle {
                                          width: 20
                                          height: 20
                                          color: selected ? "#9898B2" : "#00ffffff"
                                          border.width: 2
                                          border.color: "#9898B2"
                                          radius: 100


                                          Icon {
                                              id: checkIcon
                                              width: 12
                                              height: 12
                                              source: "qrc:/assets/check.svg"
                                              visible: selected
                                              anchors.verticalCenter: parent.verticalCenter
                                              anchors.horizontalCenter: parent.horizontalCenter
                                          }
                                      }
                                  }

                              ColumnLayout{
                                  Text {
                                      Layout.preferredWidth: mastodonServers.width - parent.x
                                      wrapMode: Text.WordWrap
                                      text: domain
                                      color: "#F3F2F7"
                                      font.weight: Font.Medium
                                  }
                                  Text {
                                      Layout.preferredWidth: mastodonServers.width - parent.x
                                      wrapMode: Text.WordWrap
                                      text: description
                                      color: "#9898B2"
                                  }
                                  Row {
                                      spacing: 12

                                      Row {
                                          spacing: 4

                                          Icon {
                                              width: 16
                                              height: 16
                                              source: "qrc:/assets/groups.svg"
                                              anchors.verticalCenter: parent.verticalCenter
                                          }

                                          Text {
                                              text: (totalUsers / 1000).toFixed(1) + "K" // could be made "smarter" for servers with less than 1k users
                                              font.weight: Font.Medium
                                              color: "#9898B2"
                                          }
                                      }

                                      Row {
                                          spacing: 4

                                          Icon {
                                              width: 16
                                              height: 16
                                              source: "qrc:/assets/translate.svg"
                                              anchors.verticalCenter: parent.verticalCenter

                                          }

                                          Text {
                                              text: language
                                              font.capitalization: Font.AllUppercase
                                              font.weight: Font.Medium
                                              color: "#9898B2"
                                          }
                                      }
                                  }
                              }
                              MouseArea {
                                 anchors.fill: serverRow
                                 onClicked: {
                                      QClient.setSelected(index)
                                  }
                              }
                          }

                          Rectangle {
                              height: 1
                              width: parent.width
                              Layout.preferredWidth: mastodonServers.width
                              color: "#9898B2"
                              visible: index < mastodonServers.count -1
                          }

                          Item { height: 1; width: 1 } // spacer
                      }
    }

    ListView {
        id: mastodonServers

        anchors.fill: parent
        anchors {
            topMargin: 12
            leftMargin: 20
            rightMargin: 20
        }

        header: Column {
            id: serversHeader
            width: parent.width
            spacing: 12
            Text {
                id: title
                text: "Mastodon is made of users on different servers."
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
                text: "Pick a server based on our intrests, region, or a general purpose one. You can still connect with everyone, regardless of server."
                anchors {
                     left: parent.left
                     right: parent.right
                }

                wrapMode: Text.WordWrap
                font.pointSize: 12
                font.weight: Font.Medium
                color: "#9898B2"
            }

            TextInput {
                id: textEdit1
                font.pointSize: 12
                text: ""
                color: "#F3F2F7"
                font.weight: Font.Medium
                width: parent.width

                topPadding: 16
                rightPadding: 24
                bottomPadding: 16
                leftPadding: 40

                layer.enabled: true
                onEditingFinished: () => {
                    console.log("Text has changed to:", text)
                    const matchingIndices = QClient.filterServers(text)
                    for( var i = 0;i < visualModel.items.count;i++ ) {
                        const element = visualModel.items.get(i);
                        const model = element.model;
                        if(model.matchingSearchTerm !== true) {
                            element.inMatchingSearchTerm = false
                            continue;
                        }
                        element.inMatchingSearchTerm = true
                    }
                }
                Icon {
                    width: 24
                    height: 24
                    source: "qrc:/assets/search.svg"
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.left: parent.left
                    anchors.leftMargin: 12
                }

                Label {
                    color: "#9898B2"
                    text: "Search server or enter URL"
                    font.weight: Font.Medium
                    font.pointSize: 12
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.left: parent.left
                    anchors.leftMargin: 40
                    visible: !textEdit1.text && !textEdit1.activeFocus
                }

                Rectangle {
                    z: -1
                    anchors.fill: parent

                    color: "#2B2937"
                    radius: 8
                    anchors.right: parent.right
                    anchors.left: parent.left
                }
            }

            Item {height: 4; width: 1} // spacer
        }

        model: visualModel
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

                Item { Layout.fillWidth: true; }

                Button {
                    color: "#F3F2F7"
                    Layout.preferredHeight: nextLabel.paintedHeight + 24
                    enabled: true
                    onClicked: {
                        QClient.getServerRules()
                        pageStack.push(Qt.resolvedUrl("ServerRules.qml"))
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
