:: ------------------------------------------------------------
:: locresファイルを生成するバッチ
:: ------------------------------------------------------------
@echo off
pushd %~dp0
setlocal

set XCOPY_CMD=C:\Windows\System32\xcopy.exe

set LOCALIZATION_CULTURE=ja,en,fr,it,de,es

:: 対象のファイル名
set TARGET_FILE=%~1
set TARGET_NAME=%~n1
if "%TARGET_NAME%"=="" (
	echo locres生成対象が指定されていません。
	goto :END
)

set TOOL_DIR=..\tools
set LOCALIZATION_DIR=..\..\Content\Localization\Game

if not exist %LOCALIZATION_DIR% (
	mkdir %LOCALIZATION_DIR%
)

:: perforce設定読み込み
call %TOOL_DIR%\LoadP4Setting.bat

:: p4のCHARSETを設定する
if %TEMP_P4AVAILABLE% == 1 (
	p4 set P4CHARSET=shiftjis
)

:: ファイルのチェックアウト
if %TEMP_P4AVAILABLE% == 1 (
	p4 -c%TEMP_P4CLIENT% -u%TEMP_P4USER% -p%TEMP_P4PORT% edit %LOCALIZATION_DIR%/%TARGET_NAME%.manifest
	p4 -c%TEMP_P4CLIENT% -u%TEMP_P4USER% -p%TEMP_P4PORT% edit %LOCALIZATION_DIR%/*/%TARGET_NAME%.archive
)

:: エクスポート
%TOOL_DIR%\LocExport.exe -o=%LOCALIZATION_DIR% -l=%LOCALIZATION_CULTURE% %TARGET_FILE% 

:: ファイルの追加
if %TEMP_P4AVAILABLE% == 1 (
	p4 -c%TEMP_P4CLIENT% -u%TEMP_P4USER% -p%TEMP_P4PORT% add %LOCALIZATION_DIR%/%TARGET_NAME%.manifest
	
	:: addにフォルダのワイルドカードが使えない
	for /R %LOCALIZATION_DIR% %%a in (*%TARGET_NAME%.archive) do (
		p4 -c%TEMP_P4CLIENT% -u%TEMP_P4USER% -p%TEMP_P4PORT% add %%~fa
	)
)

echo locresファイルを生成しています。

set INI_NAME=_MakeLocRes.ini
set INI_DIR=Intermediate\Config\Localization\
set TARGET_DIR=..\..\%INI_DIR%
set EDITOR_EXE=..\..\..\Engine\Binaries\Win64\UE4Editor.exe
:: set UE4_EXT_OPT=(外部から設定)

:: iniファイルのコピー
%XCOPY_CMD% .\%INI_NAME% %TARGET_DIR% /q /i /y /r >NUL

:: iniファイルの書き換え
type nul > %TARGET_DIR%%INI_NAME%

setlocal enabledelayedexpansion
for /f "delims=" %%I in (%INI_NAME%) do (
    set line=%%I
    (echo !line:%%s=%TARGET_NAME%!) >> %TARGET_DIR%%INI_NAME%
)
endlocal

:: locres生成
%EDITOR_EXE% TresGame -run=UnrealEd.GatherTextCommandlet -EnableSCC -DisableSCCSubmit -config=%INI_DIR%%INI_NAME% -log=MakeLocRes%TARGET_NAME%.log %UE4_EXT_OPT%

if not %ERRORLEVEL% == 0 (
	echo locresファイルの生成に失敗しました。ログファイルを確認してください。TresGame\Saved\Logs\MakeLocRes%TARGET_NAME%.log
	if not defined UE4_EXT_OPT pause
) else (
	echo locresファイルの生成が完了しました。
)

:END

endlocal
popd
