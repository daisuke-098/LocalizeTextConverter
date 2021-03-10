@echo off

if "%1"=="" (
	:: 全コンバート
	for %%f in (*.xlsx) do (
		call :CNVFILE %%f
	)
) else (
	:: 指定コンバート
	for %%f in (%*) do (
	  call _MakeLocres.bat %%f
	)
)

goto :EOF

:CNVFILE
set FILE=%1
if not "%FILE:~0,1%"=="_" (
	call _MakeLocres.bat %1
)
exit /b
