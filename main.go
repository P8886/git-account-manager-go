package main

import (
	"fmt"
	"time"

	"git-account-manager-go/internal/gitops"
	"git-account-manager-go/internal/storage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Git Account Manager")
	myWindow.Resize(fyne.NewSize(800, 600))

	// State
	config, _ := storage.LoadConfig()

	// 如果首次运行且无账户，尝试读取当前全局配置
	if len(config.Accounts) == 0 {
		name, email, _ := gitops.GetGlobalConfig()
		if name != "" && email != "" {
			newAcc := storage.Account{
				ID:    fmt.Sprintf("%d", time.Now().UnixMilli()),
				Name:  name,
				Email: email,
			}
			config.Accounts = append(config.Accounts, newAcc)
			config.ActiveID = newAcc.ID
			storage.SaveConfig(config)
		}
	}

	// UI Components
	statusLabel := widget.NewLabel("正在读取状态...")
	updateStatus := func() {
		name, email, _ := gitops.GetGlobalConfig()
		statusLabel.SetText(fmt.Sprintf("当前全局身份: %s <%s>", name, email))
	}
	updateStatus()

	accountList := widget.NewList(
		func() int {
			return len(config.Accounts)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.AccountIcon()),
				widget.NewLabel("Template Name"),
				layout.NewSpacer(),
				widget.NewLabel("Template Email"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			acc := config.Accounts[i]
			box := o.(*fyne.Container)
			nameLabel := box.Objects[1].(*widget.Label)
			emailLabel := box.Objects[3].(*widget.Label)

			nameLabel.SetText(acc.Name)
			emailLabel.SetText(acc.Email)

			if acc.ID == config.ActiveID {
				nameLabel.TextStyle = fyne.TextStyle{Bold: true}
				emailLabel.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				nameLabel.TextStyle = fyne.TextStyle{}
				emailLabel.TextStyle = fyne.TextStyle{}
			}
		},
	)

	// Detail / Action Area
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Git 用户名"
	emailEntry := widget.NewEntry()
	emailEntry.PlaceHolder = "Git 邮箱"
	sshEntry := widget.NewEntry()
	sshEntry.PlaceHolder = "SSH 私钥路径 (可选)"

	sshSelectBtn := widget.NewButtonWithIcon("浏览", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				sshEntry.SetText(reader.URI().Path())
			}
		}, myWindow)
	})

	form := container.NewVBox(
		widget.NewLabel("添加/编辑账户"),
		widget.NewForm(
			widget.NewFormItem("用户名", nameEntry),
			widget.NewFormItem("邮箱", emailEntry),
			widget.NewFormItem("SSH Key", container.NewBorder(nil, nil, nil, sshSelectBtn, sshEntry)),
		),
	)

	var refreshList func()

	addBtn := widget.NewButtonWithIcon("保存账户", theme.DocumentSaveIcon(), func() {
		if nameEntry.Text == "" || emailEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("用户名和邮箱不能为空"), myWindow)
			return
		}

		newAcc := storage.Account{
			ID:         fmt.Sprintf("%d", time.Now().UnixMilli()),
			Name:       nameEntry.Text,
			Email:      emailEntry.Text,
			SSHKeyPath: sshEntry.Text,
		}

		config.Accounts = append(config.Accounts, newAcc)
		storage.SaveConfig(config)

		nameEntry.SetText("")
		emailEntry.SetText("")
		sshEntry.SetText("")
		refreshList()
	})

	deleteBtn := widget.NewButtonWithIcon("删除选中", theme.DeleteIcon(), func() {
		// Logic handled in OnSelected
	})
	deleteBtn.Disable()

	switchBtn := widget.NewButtonWithIcon("切换到选中", theme.ConfirmIcon(), func() {
		// Logic handled in OnSelected
	})
	switchBtn.Disable()

	var currentSelectedID int = -1

	accountList.OnSelected = func(id widget.ListItemID) {
		currentSelectedID = id
		deleteBtn.Enable()
		switchBtn.Enable()

		// Fill form for quick copy/edit (optional)
		acc := config.Accounts[id]
		nameEntry.SetText(acc.Name)
		emailEntry.SetText(acc.Email)
		sshEntry.SetText(acc.SSHKeyPath)
	}

	accountList.OnUnselected = func(id widget.ListItemID) {
		currentSelectedID = -1
		deleteBtn.Disable()
		switchBtn.Disable()
	}

	deleteBtn.OnTapped = func() {
		if currentSelectedID >= 0 {
			id := currentSelectedID
			dialog.ShowConfirm("删除账户", "确定要删除吗？", func(b bool) {
				if b {
					// Remove from slice
					config.Accounts = append(config.Accounts[:id], config.Accounts[id+1:]...)
					storage.SaveConfig(config)
					accountList.UnselectAll()
					refreshList()
				}
			}, myWindow)
		}
	}

	switchBtn.OnTapped = func() {
		if currentSelectedID >= 0 {
			id := currentSelectedID
			acc := config.Accounts[id]
			err := gitops.SetGlobalConfig(acc.Name, acc.Email, acc.SSHKeyPath)
			if err != nil {
				dialog.ShowError(err, myWindow)
			} else {
				config.ActiveID = acc.ID
				storage.SaveConfig(config)
				updateStatus()
				refreshList()
				dialog.ShowInformation("成功", fmt.Sprintf("已切换到 %s", acc.Name), myWindow)
			}
		}
	}

	refreshList = func() {
		accountList.Refresh()
	}

	// Layout
	leftSide := container.NewBorder(nil, nil, nil, nil, accountList)
	rightSide := container.NewVBox(
		form,
		addBtn,
		layout.NewSpacer(),
		widget.NewSeparator(),
		container.NewHBox(switchBtn, deleteBtn),
	)

	split := container.NewHSplit(leftSide, rightSide)
	split.SetOffset(0.4)

	mainContent := container.NewBorder(
		container.NewVBox(statusLabel, widget.NewSeparator()),
		nil, nil, nil,
		split,
	)

	myWindow.SetContent(mainContent)
	myWindow.ShowAndRun()
}
