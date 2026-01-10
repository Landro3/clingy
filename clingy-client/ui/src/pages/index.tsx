import { useNavigation } from '../context/navigation';
import { Pages as PageEnum } from '../context/navigation';
import Chat from './chat';
import Config from './config';
import Contacts from './contacts';
import Intro from './intro';

export default function Pages() {
  const { currentPage } = useNavigation();

  switch (currentPage) {
    case PageEnum.Config:
      return <Config />;
    case PageEnum.Chat:
      return <Chat />;
    case PageEnum.Contacts:
      return <Contacts />;
    default:
      return <Intro />;
  }
}
