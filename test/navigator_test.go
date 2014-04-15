package test

import (
	"github.com/jmacdonald/liberator/filesystem/directory"
	"github.com/jmacdonald/liberator/view"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Navigator", func() {
	var (
		navigator    *directory.Navigator
		path         string
		error        error
		originalPath string
	)

	BeforeEach(func() {
		originalPath, _ = os.Getwd()
		navigator = directory.NewNavigator(originalPath)
	})

	Describe("NewNavigator", func() {
		It("sets the current path using its path argument", func() {
			Expect(navigator.CurrentPath()).To(Equal(originalPath))
		})
	})

	Describe("SetWorkingDirectory", func() {
		// Change the working directory right before every test.
		JustBeforeEach(func() {
			error = navigator.SetWorkingDirectory(path)
		})

		Context("path is a valid directory", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				path += "/sample"
			})

			It("returns a nil error", func() {
				Expect(error).To(BeNil())
			})

			It("updates current path with its path argument", func() {
				Expect(navigator.CurrentPath()).To(Equal(path))
			})

			It("updates entries using path argument", func() {
				Expect(navigator.Entries()).To(Equal(directory.Entries(path)))
			})

			It("resets selected index to zero", func() {
				navigator.SelectNextEntry()
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))

				navigator.SetWorkingDirectory(path)
				Expect(navigator.SelectedIndex()).To(BeZero())
			})

			It("resets previous view data indices", func() {
				_, _ = navigator.View(1)

				navigator.SetWorkingDirectory(path)
				Expect(navigator.ViewDataIndices()).To(Equal([2]uint16{0, 0}))
			})
		})

		Context("path is a file", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				path += "/sample/file"

				// Increment the selected index so we can ensure
				// it isn't reset to zero later on.
				navigator.SelectNextEntry()
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})

			It("does not update current path", func() {
				Expect(navigator.CurrentPath()).To(Equal(originalPath))
			})

			It("does not update entries", func() {
				Expect(navigator.Entries()).To(Equal(directory.Entries(originalPath)))
			})

			It("does not reset selected index to zero", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})
		})

		Context("path is invalid", func() {
			BeforeEach(func() {
				path = "/asdf"

				// Increment the selected index so we can ensure
				// it isn't reset to zero later on.
				navigator.SelectNextEntry()
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})

			It("does not update current path", func() {
				Expect(navigator.CurrentPath()).To(Equal(originalPath))
			})

			It("does not update entries", func() {
				Expect(navigator.Entries()).To(Equal(directory.Entries(originalPath)))
			})

			It("does not reset selected index to zero", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})
		})

		Context("path has a trailing slash", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				path += "/"
			})

			It("strips the trailing slash", func() {
				Expect(navigator.CurrentPath()).To(Equal(path[:len(path)-1]))
			})
		})
	})

	Describe("SelectedEntry", func() {
		Context("the second entry is selected", func() {
			BeforeEach(func() {
				navigator.SelectNextEntry()
			})

			It("returns the entry at the currently selected index", func() {
				Expect(navigator.SelectedEntry()).To(Equal(navigator.Entries()[navigator.SelectedIndex()]))
			})
		})
	})

	Describe("SelectNextEntry", func() {
		JustBeforeEach(func() {
			navigator.SelectNextEntry()
		})

		Context("directory has never been set", func() {
			BeforeEach(func() {
				navigator = new(directory.Navigator)
			})

			It("does not change the selected index", func() {
				Expect(navigator.SelectedIndex()).To(BeZero())
			})
		})

		Context("directory has been set and has entries", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				navigator.SetWorkingDirectory(path)
			})

			It("increments the selected index by one", func() {
				Expect(navigator.SelectedIndex()).To(BeEquivalentTo(1))
			})

			Context("last entry is selected", func() {
				var selectedIndex uint16

				BeforeEach(func() {
					// Call SelectNextEntry() until the last entry is selected.
					for uint16(len(navigator.Entries()))-navigator.SelectedIndex() > 1 {
						navigator.SelectNextEntry()
					}

					// Keep a reference to the last index.
					selectedIndex = navigator.SelectedIndex()
				})

				It("does not change the selected index", func() {
					Expect(navigator.SelectedIndex()).To(Equal(selectedIndex))
				})
			})
		})
	})

	Describe("SelectPreviousEntry", func() {
		JustBeforeEach(func() {
			navigator.SelectPreviousEntry()
		})

		Context("directory has never been set", func() {
			It("does not change the selected index", func() {
				Expect(navigator.SelectedIndex()).To(BeZero())
			})
		})

		Context("directory has been set and has entries", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				navigator.SetWorkingDirectory(path)
			})

			It("does not change the selected index", func() {
				Expect(navigator.SelectedIndex()).To(BeZero())
			})

			Context("last entry is selected", func() {
				var selectedIndex uint16

				BeforeEach(func() {
					// Call SelectNextEntry() until the last entry is selected.
					for uint16(len(navigator.Entries()))-navigator.SelectedIndex() > 1 {
						navigator.SelectNextEntry()
					}

					// Keep a reference to the last index.
					selectedIndex = navigator.SelectedIndex()
				})

				It("decrements the selected index by one", func() {
					Expect(navigator.SelectedIndex()).To(BeEquivalentTo(selectedIndex - 1))
				})
			})
		})
	})

	Describe("IntoSelectedEntry", func() {
		JustBeforeEach(func() {
			error = navigator.IntoSelectedEntry()
		})

		BeforeEach(func() {
			path, _ = os.Getwd()
			path += "/sample"
			navigator.SetWorkingDirectory(path)
		})

		Context("a directory is selected", func() {
			BeforeEach(func() {
				for navigator.SelectedEntry().Name != "directory" {
					navigator.SelectNextEntry()
				}
			})

			It("navigates into the selected entry", func() {
				Expect(navigator.CurrentPath()).To(BeEquivalentTo(path + "/directory"))
			})

			It("does not return an error", func() {
				Expect(error).To(BeNil())
			})
		})

		Context("a file is selected", func() {
			BeforeEach(func() {
				for navigator.SelectedEntry().Name != "file" {
					navigator.SelectNextEntry()
				}
			})

			It("does not change the working directory", func() {
				Expect(navigator.CurrentPath()).To(BeEquivalentTo(path))
			})

			It("returns an error", func() {
				Expect(error).ToNot(BeNil())
			})
		})
	})

	Describe("RemoveSelectedEntry", func() {
		var file_name, directory_name string

		JustBeforeEach(func() {
			error = navigator.RemoveSelectedEntry()
		})

		Context("selected entry is a file", func() {
			BeforeEach(func() {
				file_name = "new_file"
				os.Create(file_name)

				// Update the navigator's cached entries.
				navigator.SetWorkingDirectory(originalPath)

				for navigator.SelectedEntry().Name != file_name {
					navigator.SelectNextEntry()
				}
			})

			It("deletes the file", func() {
				_, err := os.Stat(file_name)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

			It("removes the file from the navigator's entries", func() {
				file_entry := &directory.Entry{Name: file_name, Size: 0, IsDirectory: false}
				Expect(navigator.Entries()).ToNot(ContainElement(file_entry))
			})
		})

		Context("selected entry is a directory with files", func() {
			BeforeEach(func() {
				file_name = "new_file"
				directory_name = "new_directory"
				os.Mkdir(directory_name, 0700)
				os.Create(directory_name + "/" + file_name)

				// Update the navigator's cached entries.
				navigator.SetWorkingDirectory(originalPath)

				for navigator.SelectedEntry().Name != directory_name {
					navigator.SelectNextEntry()
				}
			})

			It("deletes the directory", func() {
				_, err := os.Stat(directory_name)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
	})

	Describe("ToParentDirectory", func() {
		var parent_path string

		JustBeforeEach(func() {
			error = navigator.ToParentDirectory()
		})

		Context("directory has a parent directory", func() {
			BeforeEach(func() {
				path, _ = os.Getwd()
				parent_path = path + "/sample"
				path += "/sample/directory"
				navigator.SetWorkingDirectory(path)
			})

			It("navigates to the parent directory", func() {
				Expect(navigator.CurrentPath()).To(BeEquivalentTo(parent_path))
			})
		})
	})

	Describe("View", func() {
		var rows []view.Row
		var status string
		var maxRows uint16

		JustBeforeEach(func() {
			rows, status = navigator.View(maxRows)
		})

		It("returns the current directory path as its status", func() {
			Expect(status).To(Equal(navigator.CurrentPath()))
		})

		Context("maxRows is set to 1", func() {
			BeforeEach(func() {
				maxRows = 1
			})

			It("returns a slice with the right number of entries", func() {
				Expect(len(rows)).To(BeEquivalentTo(maxRows))
			})

			It("stores the proper view data indices", func() {
				Expect(navigator.ViewDataIndices()).To(Equal([2]uint16{0, 1}))
			})

			Describe("returned row", func() {
				It("has its left value set to the first entry's name", func() {
					Expect(rows[0].Left).To(Equal(navigator.Entries()[0].Name))
				})

				It("has its right value set to the first entry's formatted size", func() {
					formattedSize := view.Size(navigator.Entries()[0].Size)
					Expect(rows[0].Right).To(Equal(formattedSize))
				})

				It("has its highlight value set to the first entry's highlighted status", func() {
					Expect(rows[0].Highlight).To(BeTrue())
				})

				Context("selected entry is a directory", func() {
					BeforeEach(func() {
						entry := navigator.SelectedEntry()

						for !entry.IsDirectory {
							navigator.SelectNextEntry()
							entry = navigator.SelectedEntry()
						}
					})

					It("has its colour value set to true", func() {
						Expect(rows[0].Colour).To(BeTrue())
					})

					It("has a forward slash appended to its name", func() {
						Expect(rows[0].Left).To(Equal(navigator.SelectedEntry().Name + "/"))
					})
				})

				Context("selected entry is not a directory", func() {
					BeforeEach(func() {
						entry := navigator.SelectedEntry()

						for entry.IsDirectory {
							navigator.SelectNextEntry()
							entry = navigator.SelectedEntry()
						}
					})

					It("has its colour value set to false", func() {
						Expect(rows[0].Colour).To(BeFalse())
					})
				})
			})
		})

		Context("maxRows is set to 2", func() {
			BeforeEach(func() {
				maxRows = 2
			})

			It("returns a slice with the right number of entries", func() {
				Expect(len(rows)).To(BeEquivalentTo(maxRows))
			})

			It("stores the proper view data indices", func() {
				Expect(navigator.ViewDataIndices()).To(Equal([2]uint16{0, 2}))
			})

			Context("selected entry has never been changed", func() {
				It("returns the first and second rows", func() {
					Expect(rows[0].Left).To(ContainSubstring(navigator.Entries()[0].Name))
					Expect(rows[1].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("sets the first row as highlighted", func() {
					Expect(rows[0].Highlight).To(BeTrue())
				})
			})

			Context("the second entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
				})

				It("returns the first row", func() {
					Expect(rows[0].Left).To(ContainSubstring(navigator.Entries()[0].Name))
				})

				It("returns the second row", func() {
					Expect(rows[1].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("sets the second row as highlighted", func() {
					Expect(rows[1].Highlight).To(BeTrue())
				})
			})

			Context("the second entry is selected, the view is rendered, and then the third entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
					_, _ = navigator.View(maxRows)
					navigator.SelectNextEntry()
				})

				It("returns the second row", func() {
					Expect(rows[0].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("returns the third row", func() {
					Expect(rows[1].Left).To(ContainSubstring(navigator.Entries()[2].Name))
				})

				It("sets the third row as highlighted", func() {
					Expect(rows[1].Highlight).To(BeTrue())
				})
			})

			Context("the third entry is selected, the view is rendered, and then the second entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
					navigator.SelectNextEntry()
					_, _ = navigator.View(maxRows)
					navigator.SelectPreviousEntry()
				})

				It("returns the second row", func() {
					Expect(rows[0].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("returns the third row", func() {
					Expect(rows[1].Left).To(ContainSubstring(navigator.Entries()[2].Name))
				})

				It("sets the second row as highlighted", func() {
					Expect(rows[0].Highlight).To(BeTrue())
				})
			})

			Context("the fourth entry is selected, the view is rendered, and then the second entry is selected", func() {
				BeforeEach(func() {
					navigator.SelectNextEntry()
					navigator.SelectNextEntry()
					navigator.SelectNextEntry()
					_, _ = navigator.View(maxRows)
					navigator.SelectPreviousEntry()
					navigator.SelectPreviousEntry()
				})

				It("returns the second row", func() {
					Expect(rows[0].Left).To(ContainSubstring(navigator.Entries()[1].Name))
				})

				It("returns the third row", func() {
					Expect(rows[1].Left).To(ContainSubstring(navigator.Entries()[2].Name))
				})

				It("sets the second row as highlighted", func() {
					Expect(rows[0].Highlight).To(BeTrue())
				})
			})
		})

		Context("in an empty directory", func() {
			var emptyDirectoryPath string

			BeforeEach(func() {
				// Create an empty directory.
				path, _ = os.Getwd()
				fileInfo, _ := os.Stat(path + "/sample/")
				emptyDirectoryPath = path + "/sample/empty"
				_ = os.Mkdir(emptyDirectoryPath, fileInfo.Mode())

				// Navigate into the empty directory.
				navigator.SetWorkingDirectory(emptyDirectoryPath)
			})

			AfterEach(func() {
				os.Remove(emptyDirectoryPath)
			})

			It("does not panic", func() {
				Expect(func() { navigator.View(1) }).ToNot(Panic())
			})

			It("returns an empty slice of rows", func() {
				slice, _ := navigator.View(1)
				Expect(len(slice)).To(BeZero())
			})
		})
	})
})
