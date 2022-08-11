# Termiter

Termiter is a make-like program for building tasks execution

## A Termiter file

In a Termiter file there are 3 types of objects: actions, profiles and commands
The purpose of each is well defined.

A termite execution starts when the user tries to run either a command or an action of the file

- Action: Can be run once or many time during the execution and always do the same.
  - Can call other actions
- Profile: Are used by commands to know what inputs they expect from the user
  - Can extend another profile
- Command: Can only be run once and by the user itself
  - Can call actions
  - Can use a profile