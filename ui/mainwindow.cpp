#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QWidget>
#include <QStackedWidget>
#include <QString>
#include <QFileDialog>
#include <QLineEdit>
#include <QList>
#include <QDebug>
#include <QMessageBox>
#include <QDesktopWidget>
#include <QSize>
#include <QStyle>
#include <QResizeEvent>
#include <QRect>
#include <iostream>
#include <cstring>
#include <stdlib.h>
#include "libs/xmnseedwords.h"
#include "libs/xmncrypto.h"

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);

    this->setWindowIcon(QIcon(":images/icon_32.png"));

    this->dbDirPath = QString{};
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_loadAccountBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(4);
}

void MainWindow::on_createAccountBackHomeBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(0);
}

void MainWindow::on_loadAccountBackHomeBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(0);
}

void MainWindow::on_selectDirBtn_clicked()
{
    QString dir = QFileDialog::getExistingDirectory(
        this, tr("Open Directory"),
        QDir::currentPath(),
        QFileDialog::ShowDirsOnly| QFileDialog::DontResolveSymlinks
    );

    ui->dirInput->setText(dir);
    this->dbDirPath = dir;

}

void MainWindow::on_createAccountNextBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(2);
}

void MainWindow::on_seedWordsBackHomeBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(0);
}

void MainWindow::on_seedWordsSelectDirBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(1);
}

void MainWindow::on_seedWordsGenBtn_clicked()
{
    // create the list of input fields:
    QList<QLineEdit*> wordLines;
    wordLines.append(ui->seedWords1);
    wordLines.append(ui->seedWords2);
    wordLines.append(ui->seedWords3);
    wordLines.append(ui->seedWords4);
    wordLines.append(ui->seedWords5);
    wordLines.append(ui->seedWords6);
    wordLines.append(ui->seedWords7);
    wordLines.append(ui->seedWords8);
    wordLines.append(ui->seedWords9);
    wordLines.append(ui->seedWords10);
    wordLines.append(ui->seedWords11);
    wordLines.append(ui->seedWords12);

    // call the generate seed words library:
    GoString lang = {"en", 2};
    GoInt amount = wordLines.count();

    // set the words in the text fields:
    char* seedWords[12] = {};
    xGetWords(seedWords, lang, amount);
    for (int i = 0; i < this->amountSeedWords; i++) {
        wordLines.at(i)->setText(seedWords[i]);
    }
}

void MainWindow::on_seedWordsSaveBtn_clicked()
{
    // create the list of input fields:
    QList<QLineEdit*> wordLines;
    wordLines.append(ui->seedWords1);
    wordLines.append(ui->seedWords2);
    wordLines.append(ui->seedWords3);
    wordLines.append(ui->seedWords4);
    wordLines.append(ui->seedWords5);
    wordLines.append(ui->seedWords6);
    wordLines.append(ui->seedWords7);
    wordLines.append(ui->seedWords8);
    wordLines.append(ui->seedWords9);
    wordLines.append(ui->seedWords10);
    wordLines.append(ui->seedWords11);
    wordLines.append(ui->seedWords12);

    char* seedWords[12] = {};
    for (int i = 0; i < this->amountSeedWords; i++) {
        const char* original = wordLines.at(i)->text().toStdString().c_str();
        char* dest = new char[strlen(original)];
        seedWords[i] = strcpy(dest, original);
    }

    // generate the encrypted pk:
    char* encPK = xGenEncryptedPk(seedWords, this->amountSeedWords);

    // save the seed words on file:
    QString filePath = this->dbDirPath.append("/key.xmn");
    QFile file(filePath);
    if (!file.open(QIODevice::ReadWrite)) {
        QMessageBox::information(this, tr("Unable to open file.  Please choose a directory we can write in."), file.errorString());

        // set the widget to change the directory:
        ui->stInstallWidget->setCurrentIndex(1);
        return;
    }

    QTextStream stream(&file);
    stream << encPK << endl;


    // create the list of text fields:
    QList<QLineEdit*> savedSeedWords;
    savedSeedWords.append(ui->savedSeedWords1);
    savedSeedWords.append(ui->savedSeedWords2);
    savedSeedWords.append(ui->savedSeedWords3);
    savedSeedWords.append(ui->savedSeedWords4);
    savedSeedWords.append(ui->savedSeedWords5);
    savedSeedWords.append(ui->savedSeedWords6);
    savedSeedWords.append(ui->savedSeedWords7);
    savedSeedWords.append(ui->savedSeedWords8);
    savedSeedWords.append(ui->savedSeedWords9);
    savedSeedWords.append(ui->savedSeedWords10);
    savedSeedWords.append(ui->savedSeedWords11);
    savedSeedWords.append(ui->savedSeedWords12);

    // set the words:
    for (int i = 0; i < this->amountSeedWords; i++) {
        QString seedWord = seedWords[i];
        QLineEdit* wordLine = savedSeedWords.at(i);
        wordLine->setText(seedWord);
        wordLine->setReadOnly(true);

    }

    // change the panel:
    ui->stInstallWidget->setCurrentIndex(3);
}

void MainWindow::on_fileSavedLoginBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(4);
}

void MainWindow::on_loadAccountSelectFileBtn_clicked()
{
    QString filename = QFileDialog::getOpenFileName(
        this,
        tr("Open File"),
        QDir::currentPath(),
        "XMN File (*.xmn)"
    );

    // set the path in the UI:
    ui->loadAccountFilePath->setText(filename);

}

void MainWindow::on_loadAccountNextBtn_clicked()
{
    ui->stInstallWidget->setCurrentIndex(5);
}

void MainWindow::on_decryptBtn_clicked()
{
    // create the list of text fields:
    QList<QLineEdit*> decPassWords;
    decPassWords.append(ui->decPassWord1);
    decPassWords.append(ui->decPassWord2);
    decPassWords.append(ui->decPassWord3);
    decPassWords.append(ui->decPassWord4);
    decPassWords.append(ui->decPassWord5);
    decPassWords.append(ui->decPassWord6);
    decPassWords.append(ui->decPassWord7);
    decPassWords.append(ui->decPassWord8);
    decPassWords.append(ui->decPassWord9);
    decPassWords.append(ui->decPassWord10);
    decPassWords.append(ui->decPassWord11);
    decPassWords.append(ui->decPassWord12);

    // get the words:
    char* passWords[12] = {};
    for (int i = 0; i < this->amountSeedWords; i++) {
        const char* original = decPassWords.at(i)->text().toStdString().c_str();
        char* dest = new char[strlen(original)];
        passWords[i] = strcpy(dest, original);
    }

    // retrieve the private/public pair from file:
    QString filename = ui->loadAccountFilePath->text();
    QFile file(filename);
    if (!file.open(QIODevice::ReadWrite)) {
        QMessageBox::information(this, tr("Unable to open file"), file.errorString());
        return;
    }

    QString encPK;
    QTextStream stream(&file);
    stream >> encPK;

    // decrypt the pk:
    char* decPK = xDecrypt((char*) encPK.toStdString().c_str(), passWords, this->amountSeedWords);

    qDebug() << "Dec PK: " << decPK;
    if ((decPK != NULL) && (decPK[0] == '\0')) {
        QMessageBox::information(this, tr("Invalid pass words"), "The pass words you entered cannot decrypt the selected stored private key.");
        return;
    }

    // resize the window:
    QSize avSize = qApp->desktop()->availableGeometry().size();
    this->setGeometry(
        QStyle::alignedRect(
            Qt::LeftToRight,
            Qt::AlignCenter,
            QSize( avSize.width() * 0.9, avSize.height() * 0.9),
            qApp->desktop()->availableGeometry()
        )
    );

    // change the panel:
    ui->sWidget->setCurrentIndex(1);

    // set the decrypted PK:
    this->decPK = decPK;
}

void MainWindow::resizeEvent(QResizeEvent* event)
{
   QMainWindow::resizeEvent(event);

   // resize and reposition the label:
   QRect geo = ui->copyrightLbl->geometry();
   ui->copyrightLbl->setGeometry(
               geo.x(),
               this->geometry().height() - geo.height(),
               this->geometry().width(),
               geo.height()
   );

   ui->copyrightLbl->updateGeometry();
}
