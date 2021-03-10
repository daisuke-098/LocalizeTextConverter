@echo off
setlocal

set IMPORT_DIR=.\

if "%1"=="" (
	:: 全コンバート
	for /r %IMPORT_DIR% %%f in (*.xlsx) do (
		echo %%f
	)
) else (
	:: 指定コンバート
	for %%f in (%*) do (
	  call _MakeLocres.bat %%f
	)
)

endlocal
exit /b
