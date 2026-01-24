import { useNavigation } from '../context/navigation';
import { Pages as PageEnum } from '../context/navigation';
import Chat from './chat/chat';
import Intro from './intro';

export default function Pages() {
  const { currentPage } = useNavigation();

  switch (currentPage) {
    case PageEnum.Chat:
      return <Chat />;
    default:
      return <Intro />;
  }
}
