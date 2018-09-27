#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <stdlib.h>

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:

    void on_loadAccountBtn_clicked();

    void on_createAccountBackHomeBtn_clicked();

    void on_loadAccountBackHomeBtn_clicked();

    void on_selectDirBtn_clicked();

    void on_createAccountNextBtn_clicked();

    void on_seedWordsBackHomeBtn_clicked();

    void on_seedWordsSelectDirBtn_clicked();

    void on_seedWordsGenBtn_clicked();

    void on_seedWordsSaveBtn_clicked();

    void on_fileSavedLoginBtn_clicked();

    void on_loadAccountSelectFileBtn_clicked();

    void on_loadAccountNextBtn_clicked();

    void on_decryptBtn_clicked();

    void resizeEvent(QResizeEvent* event);

private:
    Ui::MainWindow *ui;
    QString dbDirPath;
    char* decPK;
    const int amountSeedWords = 12;
};

#endif // MAINWINDOW_H
