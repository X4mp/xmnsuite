#-------------------------------------------------
#
# Project created by QtCreator 2018-09-19T21:55:14
#
#-------------------------------------------------

QT       += core gui widgets

TARGET = ui
TEMPLATE = app

# The following define makes your compiler emit warnings if you use
# any feature of Qt which has been marked as deprecated (the exact warnings
# depend on your compiler). Please consult the documentation of the
# deprecated API in order to know how to port your code away from it.
DEFINES += QT_DEPRECATED_WARNINGS

# You can also make your code fail to compile if you use deprecated APIs.
# In order to do so, uncomment the following line.
# You can also select to disable deprecated APIs only up to a certain version of Qt.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

CONFIG += c++11

SOURCES += \
        main.cpp \
        mainwindow.cpp

HEADERS += \
        mainwindow.h \
        libs/xmnseedwords.h \
        libs/xmncrypto.h

FORMS += \
        mainwindow.ui

ICON = xmn.icns

CONFIG += mobility
MOBILITY =


# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target

RESOURCES += \
    images/logo.qrc

DISTFILES += \
    images/xmnservices.png \
    libs/xmncrypto.dylib \
    libs/xmnseedwords.dylib

win32:CONFIG(release, debug|release): LIBS += "$$PWD/libs/release/xmncrypto.dll"
else:win32:CONFIG(debug, debug|release): LIBS += "$$PWD/libs/debug/xmncrypto.dll"
else:macx: LIBS += "$$PWD/libs/xmncrypto.dylib"
else:unix: LIBS += "$$PWD/libs/xmncrypto.so"

win32:CONFIG(release, debug|release): LIBS += "$$PWD/libs/release/xmnseedwords.dll"
else:win32:CONFIG(debug, debug|release): LIBS += "$$PWD/libs/debug/xmnseedwords.dll"
else:macx: LIBS += "$$PWD/libs/xmnseedwords.dylib"
else:unix: LIBS += "$$PWD/libs/xmnseedwords.so"

INCLUDEPATH += $$PWD/libs
DEPENDPATH += $$PWD/libs
