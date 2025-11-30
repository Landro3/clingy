interface FocusTextBoxProps {
  text: string;
  focused: boolean;
}

export default function FocusTextBox({ text, focused }: FocusTextBoxProps) {

  return (
    <box alignItems="center" border borderStyle="rounded" borderColor={focused ? 'red' : 'white'}>
      <text fg={focused ? 'red' : 'white'}>{text}</text>
    </box>
  );
}
