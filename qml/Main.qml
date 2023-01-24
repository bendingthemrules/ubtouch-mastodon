/*
 * Copyright (C) 2022  Development@bendingtherules.nl
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 3.
 *
 * mastodon is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import QtQuick 2.7
import QtQuick.Layouts 1.3

import Ubuntu.Components 1.3
import Qt.labs.settings 1.0
import Ubuntu.PushNotifications 0.1
import QtWebEngine 1.7
import Morph.Web 0.1
import Ubuntu.Layouts 1.0

MainView {
    id: root
    objectName: 'mainView'
    applicationName: 'nl.btr.mastodon'
    automaticOrientation: true

    PushClient {
            id: pushClient
            Component.onCompleted: {
                notificationsChanged.connect((msgs) => {
                    rootObject.pushNotifications.handle(msgs[0].toString())
                });
                error.connect((err) => {
                    console.log('GOT ERROR', err);
                });
             }
            onTokenChanged: rootObject.pushNotifications.initialize(pushClient.token)
            appId: 'nl.btr.mastodon_mastodon'
    }

    Connections {
        target: UriHandler
        onOpened: {
            console.log('Open from UriHandler')

            if (uris.length > 0) {
                console.log('Clicked pushmessage while in app ' + uris[0]);
                const authCodeRegex = /mastodon:\/\/oauth\?code=(.*)/
                const authCodeMatch = uris[0].match(authCodeRegex)
                if(authCodeMatch && authCodeMatch[1]){
                    QClient.handleAuthCode(authCodeMatch[1])
                }
            }
        }
    }

    PageStack {
        id: pageStack
        Component.onCompleted: QClient.shouldSkipSelection() ? push(Qt.resolvedUrl("Webview.qml")) : push(Qt.resolvedUrl("StartScreen.qml"))
    }
}